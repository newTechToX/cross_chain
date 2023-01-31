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
	stmt := fmt.Sprintf("select * from %s where direction='out' and (chain='ethereum' or chain='bsc') and from_address_error is null limit 50000", table)
	var datas []*model.Data
	var err error
	err = a.svc.ProjectsDao.DB().Select(&datas, stmt)
	if err != nil {
		return err
	}

	size := 1000
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
		if data.ToAddress == data.Token {
			tag.IsFakeToken = 1
		}

		tx, err := a.replayer.Replay(data)
		if err != nil || tx == nil {
			err = fmt.Errorf("failed to replay %s, error: %s", data.Hash, err)
			continue
		}

		//检查from_address的行为
		amount := data.Amount.String()
		realToken := a.getRealToken(data.Token, tx.BalanceChanges)
		if len(realToken) == 0 {
			tag.TokenProfitError = 1
		}
		for _, e := range tx.BalanceChanges {
			if e.Account == data.FromAddress {
				if !a.checkFrom(realToken, amount, e) {
					println("from_addr_error ", data.Hash)
					tag.FromAddressError = 1
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

			/*if tx.BalanceChanges[i].Account == data.Token {
				for _, a := range tx.BalanceChanges[i].Assets {
					if a.Address == data.Token && a.Amount == amount {
						token_profit_ok = true
						break
					}
				}
			}*/
		}
		if tag.TokenProfitError != 1 {
			if ok := a.checkTokenProfit(realToken, tx.BalanceChanges); !ok {
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
func (a *Logic) getRealToken(token string, balanceChanges []*replay.SimAccountBalance) map[string]*big.Float {
	underlying := make(map[string]*big.Float)
	burn_address := "0x0000000000000000000000000000000000000000"

	for _, e := range balanceChanges {
		if e.Account == token || e.Account == burn_address {
			for _, ee := range e.Assets {
				if ee.Amount[0] != '-' {
					underlying[ee.Address] = a.replayer.GetAmount(ee)
				}
			}
			break
		}
	}
	return underlying
}

// 1. from转出的token与token_address收到的token是否一致
// 2. from的value之和一定小于0
func (a *Logic) checkFrom(underlying map[string]*big.Float, amount string, balance *replay.SimAccountBalance) bool {
	ok := false
	for _, asset := range balance.Assets {
		if _, exsit := underlying[asset.Address]; exsit && asset.Amount[1:] == amount {
			ok = true
		}
	}
	value := a.replayer.CalAmount(balance)
	if value.String() >= "0" {
		ok = false
	}

	return ok
}

// 返回所有获利的account
func (a *Logic) profitAccounts(balances []*replay.SimAccountBalance) map[int]*big.Float {
	accounts := make(map[int]*big.Float)
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

	stmt := fmt.Sprintf("update %s set %s, profit=$1 where id = %d", table, logicUpdateTags, Id)
	// fmt.Println(stmt)
	_, err = a.svc.ProjectsDao.DB().Exec(stmt, profit, tag.FromAddressError, tag.ToAddressProfit, tag.TokenProfitError)
	if err != nil {
		log.Error("update failed", "err", err)
	}
	return err
}

// 应用场景：有多次cross out至token_addr，但是replay的结果只有最终的value
// 需要对token_addr所有assest进行检查，计算从其他账户转来的amount之和是否等于token profit amount
func (a *Logic) checkTokenProfit(realToken map[string]*big.Float, balanceChange []*replay.SimAccountBalance) bool {
	//检查token地址实际获利的所有assest
	for token_addr, value := range realToken {
		sum, _ := new(big.Float).SetPrec(uint(256)).SetString("0")
		for _, e := range balanceChange {
			for _, ee := range e.Assets {
				if ee.Address == token_addr && ee.Amount[0] == '-' {
					v := a.replayer.GetAmount(ee)
					sum.Add(sum, v)
					break
				}
			}
		}
		sumBuf, _ := sum.MarshalText()
		valueBuf, _ := value.MarshalText()
		length := len(string(valueBuf))
		var r = string(sumBuf)[1:]
		var x = r
		var y = string(valueBuf)
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
