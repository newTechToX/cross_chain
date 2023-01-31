package anyswap

import (
	"app/model"
	"app/svc"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var _ model.EventCollector = &Anyswap{}

type Anyswap struct {
	svc *svc.ServiceContext
}

func NewAnyswapCollector(svc *svc.ServiceContext) *Anyswap {
	return &Anyswap{
		svc: svc,
	}
}

func (a *Anyswap) Name() string {
	return "Anyswap"
}

func (a *Anyswap) Contracts(chain string) []string {
	if _, ok := AnyswapContracts[chain]; !ok {
		return nil
	}
	return AnyswapContracts[chain]
}

func (a *Anyswap) Topics0(chain string) []string {
	return []string{LogAnySwapIn, LogAnySwapOut}
}

func (a *Anyswap) Extract(chain string, events model.Events) model.Results {
	ret := make(model.Results, 0)
	for _, e := range events {
		res := &model.Result{
			Chain:    chain,
			Number:   e.Number,
			Ts:       e.Ts,
			Index:    e.Index,
			Hash:     e.Hash,
			ActionId: e.Id,
			Project:  a.Name(),
			Contract: e.Address,
		}
		switch e.Topics[0] {
		case LogAnySwapIn:
			if len(e.Topics) < 4 || len(e.Data) < 2+3*64 {
				continue
			}
			res.Direction = model.InDirection
			fromChainId, _ := new(big.Int).SetString(e.Data[2+64:2+128], 16)
			res.FromChainId = (*model.BigInt)(fromChainId)
			toChainId, _ := new(big.Int).SetString(e.Data[2+128:2+192], 16)
			res.ToChainId = (*model.BigInt)(toChainId)
			res.ToAddress = "0x" + e.Topics[3][26:]
			res.Token = "0x" + e.Topics[2][26:]
			amount, _ := new(big.Int).SetString(e.Data[2:2+64], 16)
			res.Amount = (*model.BigInt)(amount)
			d := &Detail{
				SrcTxHash: e.Topics[1],
			}
			detail, err := json.Marshal(d)
			if err == nil {
				res.Detail = detail
			}
			res.MatchTag = e.Topics[1]

		case LogAnySwapOut:
			if len(e.Topics) < 4 || len(e.Data) < 2+3*64 {
				continue
			}
			res.Direction = model.OutDirection
			fromChainId, _ := new(big.Int).SetString(e.Data[2+64:2+128], 16)
			res.FromChainId = (*model.BigInt)(fromChainId)
			res.FromAddress = "0x" + e.Topics[2][26:]
			toChainId, _ := new(big.Int).SetString(e.Data[2+128:2+192], 16)
			res.ToChainId = (*model.BigInt)(toChainId)
			res.ToAddress = "0x" + e.Topics[3][26:]
			res.Token = "0x" + e.Topics[1][26:]
			amount, _ := new(big.Int).SetString(e.Data[2:2+64], 16)
			res.Amount = (*model.BigInt)(amount)
			res.MatchTag = e.Hash
		}
		ret = append(ret, res)
	}
	return ret
}

func (a *Anyswap) GetUnderlying(chain, anyToken string) (string, error) {
	p := a.svc.Providers.Get(chain)
	if p == nil {
		return "", fmt.Errorf("providers does not support %v", chain)
	}
	raw, err := p.ContinueCall("", anyToken, Underlying, nil, nil)
	if err != nil {
		return "", err
	}
	return strings.ToLower(common.BytesToAddress(raw).Hex()), nil
}
