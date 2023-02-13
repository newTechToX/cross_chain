package renbridge

import (
	"app/model"
	"app/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

type RenBridge struct {
}

var _ model.MsgCollector = &RenBridge{}

func NewRenbridgeCollector() *RenBridge {
	return &RenBridge{}
}

func (r *RenBridge) Name() string {
	return "RenBridge"
}

func (r *RenBridge) Contracts(chain string) []string {
	if _, ok := contracts[chain]; !ok {
		return nil
	}
	addrs := make([]string, 0)
	for addr := range contracts[chain] {
		addrs = append(addrs, addr)
	}
	return addrs
}

func (r *RenBridge) Selectors(chain string) []string {
	return []string{Burn, Mint}
}

func (r *RenBridge) Extract(chain string, msgs []*model.Call) model.Datas {
	if _, ok := contracts[chain]; !ok {
		return nil
	}
	ret := make(model.Datas, 0)
	for _, msg := range msgs {
		if _, ok := contracts[chain][msg.To]; !ok {
			continue
		}
		if len(msg.Input) <= 10 {
			continue
		}
		sig, rawParam := msg.Input[:10], msg.Input[10:]
		params, err := Decode(sig, rawParam)
		if err != nil {
			log.Error("decode ren bridge failed", "chain", chain, "hash", msg.Hash, "err", err)
			continue
		}
		res := &model.Data{
			Chain:    chain,
			Number:   msg.Number,
			Index:    msg.Index,
			Hash:     msg.Hash,
			ActionId: msg.Id,
			Contract: msg.To,
			// non common
			Token: contracts[chain][msg.To].Token,
		}
		switch sig {
		case Burn:
			if len(params) < 2 {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.Direction = model.OutDirection
			res.FromChainId = (*model.BigInt)(utils.GetChainId(chain))
			res.FromAddress = msg.From
			res.ToChainId = (*model.BigInt)(new(big.Int).Set(contracts[chain][msg.To].ChainId))
			to, ok := params[0].([]byte)
			if !ok {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.ToAddress = hexutil.Encode(to)
			amount, ok := params[1].(*big.Int)
			if !ok {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.Amount = (*model.BigInt)(amount)
		case Mint:
			if len(params) < 4 {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.Direction = model.InDirection
			res.FromChainId = (*model.BigInt)(new(big.Int).Set(contracts[chain][msg.To].ChainId))
			res.ToChainId = (*model.BigInt)(utils.GetChainId(chain))
			res.ToAddress = msg.From
			amount, ok := params[1].(*big.Int)
			if !ok {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.Amount = (*model.BigInt)(amount)
		default:
			continue
		}
		ret = append(ret, res)
	}
	return ret
}
