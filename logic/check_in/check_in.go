package check_in

import (
	crosschain "app/cross_chain"
	"app/dao"
	"app/logic/aml"
	"app/logic/replay"
	"app/model"
	"app/svc"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
)

type InChecker struct {
	dao      *dao.Dao
	replayer *replay.Replayer
}

type Tags struct {
	ToAddressProfit  int
	TokenProfitError int
}

const (
	SAFE                        = 0
	PROJECT_TRANSFER_TYPE_ERROR = 1
	PROJECT_TRANSFER_MORE       = 2
	PROJECT_TRANSFER_MINUS      = 3
)

func NewInChecker(svc *svc.ServiceContext) *InChecker {
	return &InChecker{
		dao:      svc.Dao,
		replayer: replay.NewReplayer(svc, aml.NewAML("../txt_config.yaml"), "../txt_config.yaml"),
	}
}

//规则1：检查是否有重复
//规则2：定时将unmatch的

func (a *InChecker) HasDuplicates(project string, datas model.Datas) map[uint64][]uint64 {
	var dup_map = make(map[uint64][]uint64)
	for _, data := range datas {
		stmt := fmt.Sprintf("select %s from %s where direction = 'in' and chain = '%s' and to_address = '%s' and from_chain = %s and to_chain = %s and id != %d and amount = %s",
			model.ResultRows, project, data.Chain, data.ToAddress, data.FromChainId.String(), data.ToChainId.String(), data.Id, data.Amount.String())
		var dups model.Datas
		err := a.dao.DB().Select(&dups, stmt)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, d := range dups {
				dup_map[data.Id] = append(dup_map[data.Id], d.Id)
			}
		}
	}
	return dup_map
}

func (a *InChecker) ReplayInTxLogic(project string, data *model.Data) (tag Tags, err error) {
	//to = token
	tag = Tags{0, 0}
	tx, err := a.replayer.Replay(data)
	if err != nil || tx == nil {
		err = fmt.Errorf("failed to replay %s, error: %s", data.Hash, err)
		return
	}

	//检查project转出的行为
	amount := data.Amount.String()
	if _, ok := crosschain.TokenTransferDirectly[project]; !ok {
		//如果不是token地址直接执行转账的项目，首先转换成map
		tag.TokenProfitError = a.checkProjectToken(project, data.Chain, data.Token, amount, tx.BalanceChanges)
	} else {

	}
	//如果数量不对，那么可能是因为一笔tx里面多笔in
	if tag.TokenProfitError == PROJECT_TRANSFER_MORE || tag.TokenProfitError == PROJECT_TRANSFER_MINUS {
		var txs model.Datas
		stmt := fmt.Sprintf("select %s from %s where hash='%s' and log_index!=%d", model.ResultRows, project, data.Hash, data.LogIndex)
		if err := a.dao.DB().Select(&txs, stmt); err != nil {
			log.Error("ReplayInTxLogic() ", "err", err)
		}
		if len(txs) != 0 {
			for _, tx := range txs {
				if tx.TokenProfitError.Int64 == SAFE {
					//如果之前已经检查过同一笔hash里面的data，直接更新为SAFE即可
					tag.TokenProfitError = SAFE
					break
				}
			}

			//说明该tx没有检查过
			if tag.TokenProfitError != SAFE {
				asset_map := a.replayer.ConvertBalanceChange2TokenMap(tx.BalanceChanges)
				total_amount_ := a.replayer.CalTokenTotalAmount(data.Token, asset_map[data.Token])
				total_amount := new(model.BigInt).SetString(total_amount_["-"].Amount, 10)
				n := new(model.BigInt).SetString("0", 10)
				for _, tx := range txs {
					n = n.Add(n, tx.Amount)
				}
				if total_amount.Cmp(n) == 0 {
					tag.TokenProfitError = SAFE
				}
			}
		}
	}

	return
}

func (a *InChecker) checkProjectToken(project, chain, token, amount string, BalanceChanges []*replay.SimAccountBalance) int {
	asset_map := a.replayer.ConvertBalanceChange2TokenMap(BalanceChanges)
	if _, ok := asset_map[token]; !ok {
		return PROJECT_TRANSFER_TYPE_ERROR
	}
	sum := new(model.BigInt).SetString("0", 10)
	for _, e := range asset_map[token] {
		if e.Amount[0] != '-' {
			continue
		}
		x := new(model.BigInt).SetString(e.Amount, 10)
		sum = sum.Add(sum, x)
	}

	y := new(model.BigInt).SetString(amount, 10)
	if sum.Cmp(y) == 0 {
		return SAFE
	} else if sum.Cmp(y) > 0 {
		return PROJECT_TRANSFER_MORE
	} else {
		return PROJECT_TRANSFER_MINUS
	}
}
