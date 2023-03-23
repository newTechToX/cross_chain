package across

import (
	"app/model"
	"math/big"
)

var _ model.EventCollector = &Across{}

type Across struct {
}

func NewAcrossCollector() *Across {
	return &Across{}
}

func (a *Across) Name() string {
	return "Across"
}

func (a *Across) Contracts(chain string) []string {
	if _, ok := AcrossContracts[chain]; !ok {
		return nil
	}
	return AcrossContracts[chain]
}

func (a *Across) Topics0(chain string) []string {
	return []string{FundsDeposited, FilledRelay}
}

func (a *Across) Extract(chain string, events model.Events) model.Datas {
	ret := make(model.Datas, 0)
	len_out := 386
	len_in := 834

	for _, e := range events {
		res := &model.Data{
			Chain:    chain,
			Number:   e.Number,
			TxIndex:  e.Index,
			Hash:     e.Hash,
			LogIndex: e.Id,
			Contract: e.Address,
			Project:  a.Name(),
		}

		switch e.Topics[0] {
		case FundsDeposited:
			if len(e.Topics) < 4 || len(e.Data) < len_out {
				continue
			}
			res.Direction = model.OutDirection
			fromChainId, _ := new(big.Int).SetString(e.Data[2+64:2+128], 16)
			res.FromChainId = (*model.BigInt)(fromChainId)
			toChainId, _ := new(big.Int).SetString(e.Data[2+128:2+192], 16)
			res.ToChainId = (*model.BigInt)(toChainId)
			res.ToAddress = "0x" + e.Data[len_out-64+24:]
			res.FromAddress = "0x" + e.Topics[3][26:]
			res.Token = "0x" + e.Topics[2][26:]
			amount, _ := new(big.Int).SetString(e.Data[2:2+64], 16)
			res.Amount = (*model.BigInt)(amount)

			depositId, _ := new(big.Int).SetString(e.Topics[1][2:], 16)
			d := &Detail{
				DepositId: depositId.String(),
			}
			res.MatchTag = d.DepositId

		case FilledRelay:
			if len(e.Topics) < 3 || len(e.Data) < len_in {
				continue
			}
			res.Direction = model.InDirection
			relayer := "0x" + e.Topics[1][26:]
			fromChainId, _ := new(big.Int).SetString(e.Data[2+64*4:2+64*5], 16)
			res.FromChainId = (*model.BigInt)(fromChainId)
			res.FromAddress = "0x" + e.Topics[2][26:]

			toChainId, _ := new(big.Int).SetString(e.Data[2+64*5:2+64*6], 16)
			res.ToChainId = (*model.BigInt)(toChainId)
			depositId, _ := new(big.Int).SetString(e.Data[len_in-64*4:len_in-64*3], 16)
			res.Token = "0x" + e.Data[len_in-64*3+24:len_in-128]
			res.ToAddress = "0x" + e.Data[len_in-64*2+24:len_in-64]
			amount, _ := new(big.Int).SetString(e.Data[2:2+64], 16)
			res.Amount = (*model.BigInt)(amount)
			d := &Detail{
				DepositId: depositId.String(),
				Relayer:   relayer,
			}
			res.MatchTag = d.DepositId
		}
		ret = append(ret, res)
	}
	return ret
}
