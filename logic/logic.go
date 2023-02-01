package logic

import (
	"app/dao"
	"app/logic/replay"
	"app/model"
	"app/provider"
	"app/svc"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"math/big"
	"strings"
)

type Logic struct {
	svc      *svc.ServiceContext
	provider *provider.Provider
	replayer *replay.Replayer
}

type Tags struct {
	FromAddressError int
	ToAddressProfit  int
	TokenProfitError int
	IsFakeToken      int
}

const from_token_amount, from_token_type, from_transfer_error = 1, 2, 3

func NewLogic(svc *svc.ServiceContext, chain string) *Logic {
	p := svc.Providers.Get(chain)
	r := (&replay.Replayer{}).NewReplayer()
	if p == nil {
		panic(fmt.Sprintf("%v: invalid provider", chain))
	}
	return &Logic{
		svc:      svc,
		provider: p,
		replayer: r,
	}
}

func (a *Logic) ReplayOutTxLogic(table string) error {
	stmt := fmt.Sprintf("select * from %s where direction='out' and isfaketoken is null and (chain='ethereum' or chain='bsc') and from_address_error is not null", table)
	var datas []*model.Data
	var err error
	err = a.svc.ProjectsDao.DB().Select(&datas, stmt)
	if err != nil {
		return err
	}
	println(len(datas))

	size := 30
	i := 0
	for i = 0; i+size < len(datas); i = i + size {
		go a.replayOutTxLogic(table, datas[i:i+size])
	}
	a.replayOutTxLogic(table, datas[i:])
	return err
}

/*
已经标注的fake token --> 做验证
*/
func (a *Logic) replayOutTxLogic(table string, datas []*model.Data) (err error) {
	for i, data := range datas {
		//to = token
		tag := Tags{0, 0, 0, 0}

		var profit []*replay.SimAccountBalance
		if data.ToAddress == data.Token && !data.IsFakeToken.Valid {
			tag.IsFakeToken = 2
		}

		tx, err := a.replayer.Replay(data)
		if err != nil || tx == nil {
			err = fmt.Errorf("failed to replay %s, error: %s", data.Hash, err)
			continue
		}

		//检查from_address的行为
		amount := data.Amount.String()
		real_token := a.getRealToken(data.Token, tx.BalanceChanges)
		if len(real_token) == 0 {
			tag.TokenProfitError = 1
		}

		tag.FromAddressError = from_transfer_error //先初始化，防止根本没有from资金动态

		for _, e := range tx.BalanceChanges {
			if e.Account == data.FromAddress {
				tag.FromAddressError = a.checkFrom(real_token, data.Token, amount, e)

				//如果arg_token != deposit_token
				if tag.FromAddressError == from_token_type {
					tag.FromAddressError = 0
					previous_token := make(map[string]*model.BigFloat)

					//获取所有arg_token的资金来源
					for token := range real_token {
						p := a.getPreviousToken(token, tx.BalanceChanges)
						for k, v := range p {
							previous_token[k] = v
						}
					}

					//检查两跳
					if flag, problem_tokens := a.checkFrom_Token(previous_token, data.Token, amount, e); flag != from_token_amount && len(problem_tokens) != 0 {
						d := &dao.Dao{}
						if a.svc == nil {
							d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
						} else {
							d = a.svc.Dao
						}

						//如果两跳仍然无法对应arg_token和deposit_token，那么就查标签
						for _, token := range problem_tokens {
							if safe, _ := d.QueryLabel(token); !safe {
								//如果deposit_token unsafe
								tag.FromAddressError += from_token_type
							} else {
								tag.FromAddressError += 0
							}
						}
					} else {
						tag.FromAddressError = flag
					}
				}
				if tag.FromAddressError != 0 {
					println("from_address_error: ", data.Hash)
				}
				break
			}
		}

		//to_address不应该获利
		idx := a.profitAccounts(tx.BalanceChanges)
		for i := range idx {
			if tx.BalanceChanges[i].Account == data.ToAddress && data.Token != data.ToAddress {
				println("toaddr_profit ", data.Hash)
				tag.ToAddressProfit = 1
			}

			if tx.BalanceChanges[i].Account != data.ToAddress {
				var d *dao.Dao
				if a.svc == nil {
					d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
				} else {
					d = a.svc.Dao
				}
				if safe, _ := d.QueryLabel(tx.BalanceChanges[i].Account); !safe {
					profit = append(profit, tx.BalanceChanges[i])
				}
			}
		}
		if tag.TokenProfitError != 1 {
			if ok := a.checkTokenProfit(real_token, tx.BalanceChanges); !ok {
				println("token_profit ", data.Hash)
				tag.TokenProfitError = 1
			}

		}
		p, _ := json.Marshal(profit)
		if len(profit) == 0 {
			p = []byte(`{}`)
		}
		err = a.logicUpdate(data.Id, p, table, tag)
		if err != nil {
			fmt.Println(err)
		}
		//from_address里面检查token是否存在，如果存在那么检查amount是否对应
		/*if changes.Account == data.FromAddress {
			ok := false
			for _, asset := range changes.Assets {
				if asset.Address == data.Token && asset.Amount[1:] == amount {
					ok = true
					break
				}
			}
			OK = ok
		}

		//token检查amount是否正确增加
		if changes.Account == data.Token {
			ok := false
			for _, asset := range changes.Assets {
				if asset.Address == data.Token && asset.Amount == amount {
					ok = true
					break
				}
			}
			OK = ok
		}

		//to_address在in里面不应该收到钱
		if changes.Account == data.ToAddress {
			ok := false
			for _, asset := range changes.Assets {
				if asset.Amount > "0" {
					ok = false
					break
				}
			}
			OK = ok
		}*/

		if i%30 == 0 && i != 0 {
			println("done: ", i)
		}
	}
	println("all done")
	return
}

// 获得实际交易的token——address
// 转的就是underlying token，直接burn掉
// 转的是wapper，溯源最终的underlying token
func (a *Logic) getRealToken(token string, balanceChanges []*replay.SimAccountBalance) map[string]*model.BigFloat {
	underlying := make(map[string]*model.BigFloat)
	burn_address := "0x0000000000000000000000000000000000000000"

	for _, e := range balanceChanges {
		//记录下token_addr收到的所有资金，也就是一跳
		if e.Account == token {
			for _, ee := range e.Assets {
				if ee.Amount[0] != '-' {
					underlying[ee.Address] = a.replayer.GetAmount(ee)
				}
			}
			break
		} else if e.Account == burn_address {
			//从burn的地址里面找是否直接burn掉token
			for _, ee := range e.Assets {
				if ee.Address == token && ee.Amount[0] != '-' {
					underlying[ee.Address] = a.replayer.GetAmount(ee)
					break
				}
			}
			if _, ok := underlying[token]; ok {
				break
			}
		}
	}
	return underlying
}

// eg: ETH -> WETH -> anyETH，该函数可以用于查找前一个token的地址
func (a *Logic) getPreviousToken(token string, balance_changes []*replay.SimAccountBalance) (previous_token map[string]*model.BigFloat) {
	previous_token = make(map[string]*model.BigFloat)
	for _, e := range balance_changes {
		for _, ee := range e.Assets {
			if ee.Address == token && ee.Amount[0] == '-' { //查找是哪个地址转出的token，获取该地址的资金来源
				for _, x := range e.Assets {
					if x.Address != token && x.Amount[0] != '-' {
						previous_token[x.Address] = a.replayer.GetAmount(x)
					}
				}
			}
		}
	}
	return
}

func (a *Logic) checkFrom(underlying map[string]*model.BigFloat, token, amount string, balance *replay.SimAccountBalance) int {

	res := 0

	//只能检查第一跳，没有查到的token暂时不需要查标签库
	if flag, problem_tokens := a.checkFrom_Token(underlying, token, amount, balance); flag == from_token_amount {
		res = from_token_amount
	} else if len(problem_tokens) != 0 {
		res = from_token_type
	}

	//depositor的所有资金之和应当小于0
	value := a.replayer.CalAmount(balance)
	if value.String() >= "0" {
		res = from_transfer_error
	}
	return res
}

// 1. from转出的token与token_address收到的token是否一致
// 2. 返回的是from转出的token中，所有有问题的token_address（amount不对或者在underlying中查不到）
func (a *Logic) checkFrom_Token(underlying map[string]*model.BigFloat, token, amount string, balance *replay.SimAccountBalance) (flag int, res []string) {
	flag = 0
	for _, asset := range balance.Assets {
		if v, ok := underlying[asset.Address]; ok {
			//如果转出的token已经被记录，或者转出的就是arg_token，那么就检查amount
			//如果amount不正确，那么该交易有问题
			x := a.replayer.GetAmount(asset)
			if x.String() != "-"+v.String() {
				flag = from_token_amount
			}
		} else if asset.Address == token && asset.Amount[1:] != amount {
			flag = from_token_amount
		} else {
			//如果deposit_token仍然查不到，那么返回token地址，查标签库
			res = append(res, asset.Address)
		}
	}
	return
}

// 返回所有获利的account
func (a *Logic) profitAccounts(balances []*replay.SimAccountBalance) map[int]*model.BigFloat {
	accounts := make(map[int]*model.BigFloat)
	for i, balance := range balances {
		value := a.replayer.CalAmount(balance)
		if value.String() > "0" {
			accounts[i] = value
		}
	}
	return accounts
}

/*
balance change的逻辑 --> 检测所有cross=out（无论是否match）
*/

// 更新表字段
var (
	logicUpdateNames = []string{"from_address_error", "to_address_profit", "token_profit_error"}
	logicUpdateRows  = strings.Join(logicUpdateNames, ",")
	logicUpdateTags  = builder.PostgreSqlJoin(logicUpdateNames)
)

func (a *Logic) logicUpdate(Id uint64, profit []byte, table string, tag Tags) error {
	var err error
	d := &dao.Dao{}
	if a.svc == nil {
		d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	} else {
		d = a.svc.ProjectsDao
	}

	stmt := fmt.Sprintf("update %s set %s, profit=$1 where id = %d", table, logicUpdateTags, Id)
	// fmt.Println(stmt)
	_, err = d.DB().Exec(stmt, profit, tag.FromAddressError, tag.ToAddressProfit, tag.TokenProfitError)
	if err != nil {
		log.Error("update failed", "err", err)
	}
	return err
}

// 应用场景：有多次cross out至token_addr，但是replay的结果只有最终的value
// 需要对token_addr所有assest进行检查，计算从其他账户转来的amount之和是否等于token profit amount
func (a *Logic) checkTokenProfit(realToken map[string]*model.BigFloat, balanceChange []*replay.SimAccountBalance) bool {
	//检查token地址实际获利的所有assest
	for token_addr, value := range realToken {
		sum, _ := new(big.Float).SetPrec(uint(256)).SetString("0")
		for _, e := range balanceChange {
			for _, ee := range e.Assets {
				if ee.Address == token_addr && ee.Amount[0] == '-' {
					//无论其他地址转入或转出该token，全都加起来
					v := a.replayer.GetAmount(ee)
					sum.Add(sum, (*big.Float)(v))
					break
				}
			}
		}

		var r = (*model.BigFloat)(sum).String()[1:]
		var y = (value).String()
		var x, length = r, len(y)
		if len(r) > length {
			if r[length] >= '5' {
				x = x[:length-1]
				y = y[:length-1]
			} else {
				x = x[:length]
			}
		}

		if x != y {
			return false
		}
	}
	return true
}
