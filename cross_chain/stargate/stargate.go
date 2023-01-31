package stargate

import (
	"app/model"
	"app/svc"
	"app/utils"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

type OutDetail struct {
	// for swap
	DstPoolId uint64 `json:"dstPoolId"`
	// for bridge
	Address string `json:"address"`
	MsgType uint64 `json:"msgType"`
	Nonce   uint64 `json:"nonce"`
}

type InDetail struct {
	SrcAddress string `json:"srcAddress"`
	DstAddress string `json:"dstAddress"`
	Nonce      uint64 `json:"nonce"`
}

var _ model.EventCollector = &Stargate{}

type Stargate struct {
	svc *svc.ServiceContext
}

func NewStargateCollector(svc *svc.ServiceContext) *Stargate {
	return &Stargate{
		svc: svc,
	}
}

func (a *Stargate) Name() string {
	return "Stargate"
}

func (a *Stargate) ChangeID(id *big.Int) *big.Int {
	num, _ := strconv.Atoi(id.String())
	switch num {
	case 101:
		num = 1
	case 102:
		num = 56
	case 106:
		num = 43114
	case 109:
		num = 137
	case 110:
		num = 42161
	case 111:
		num = 10
	case 112:
		num = 250
	}

	id = new(big.Int).SetInt64(int64(num))

	return id
}

func (a *Stargate) Contracts(chain string) []string {
	if _, ok := StargateContracts[chain]; !ok {
		return nil
	}
	return StargateContracts[chain]
}

func (a *Stargate) Topics0(chain string) []string {
	return []string{Swap, SendMsg, RedeemLocalCallback, RedeemRemote,
		SendToChain, SwapRemote, PacketReceived, ReceiveFromChain}
}

func (a *Stargate) Extract(chain string, events model.Events) model.Results {
	if _, ok := StargateContracts[chain]; !ok {
		return nil
	}
	ret := make(model.Results, 0)

	outPairs := FindParis(events, Swap, SendMsg)
	for _, outPair := range outPairs {
		if len(outPair[0].Topics) != 1 || len(outPair[0].Data) < 2+8*64 {
			continue
		}
		if len(outPair[1].Topics) != 1 || len(outPair[1].Data) < 2+2*64 {
			continue
		}
		res := &model.Result{
			Chain:       chain,
			Number:      outPair[0].Number,
			Ts:          outPair[0].Ts,
			Index:       outPair[0].Index,
			Hash:        outPair[0].Hash,
			ActionId:    outPair[0].Id,
			Project:     a.Name(),
			Contract:    outPair[0].Address,
			Direction:   model.OutDirection,
			FromChainId: (*model.BigInt)(utils.GetChainId(chain)),
		}
		from := "0x" + outPair[0].Data[2+64*2+24:2+64*3]
		// if !utils.Contains(from, StargateContracts[chain]) {
		res.FromAddress = from
		// }
		toChainId, _ := new(big.Int).SetString(outPair[0].Data[2:66], 16)
		toChainId = a.ChangeID(toChainId)
		res.ToChainId = (*model.BigInt)(toChainId)
		token, err := a.GetPoolToken(chain, outPair[0].Address)
		if err == nil {
			res.Token = token
		} else {
			log.Error("stargate: cannot get pool token", "chain", chain, "hash", outPair[0].Hash, "pool", outPair[0].Address, "err", err)
			continue
		}
		amount := utils.HexSum(outPair[0].Data[2+64*3:2+64*4], outPair[0].Data[2+64*4:2+64*5],
			outPair[0].Data[2+64*5:2+64*6], outPair[0].Data[2+64*6:2+64*7], outPair[0].Data[2+64*7:2+64*8])
		convRate, err := a.GetPoolConvertRate(chain, outPair[0].Address)
		if err != nil {
			log.Error("stargate: cannot get pool convert rate", "chain", chain, "hash", outPair[0].Hash, "pool", outPair[0].Address, "err", err)
			continue
		}
		amount.Mul(amount, convRate)
		res.Amount = (*model.BigInt)(amount)

		detail := &OutDetail{
			DstPoolId: utils.ParseStrToUint64("0x" + outPair[0].Data[2+64:2+64*2]),
			Address:   outPair[1].Address,
			MsgType:   utils.ParseStrToUint64(outPair[1].Data[:66]),
			Nonce:     utils.ParseStrToUint64("0x" + outPair[1].Data[66:]),
		}
		res.Detail, _ = json.Marshal(detail)
		nonce, _ := new(big.Int).SetString(outPair[1].Data[66:], 16)
		res.MatchTag = nonce.String()
		ret = append(ret, res)
	}

	inPairs := FindParis(events, PacketReceived, SwapRemote)
	for _, inPair := range inPairs {
		if len(inPair[0].Topics) != 3 || len(inPair[0].Data) < 2+3*64 {
			continue
		}
		if len(inPair[1].Topics) != 1 || len(inPair[1].Data) < 2+4*64 {
			continue
		}
		res := &model.Result{
			Chain:     chain,
			Number:    inPair[1].Number,
			Ts:        inPair[1].Ts,
			Index:     inPair[1].Index,
			Hash:      inPair[1].Hash,
			ActionId:  inPair[1].Id,
			Project:   a.Name(),
			Contract:  inPair[1].Address,
			Direction: model.InDirection,
			ToChainId: (*model.BigInt)(utils.GetChainId(chain)),
		}
		fromChainId, _ := new(big.Int).SetString(inPair[0].Topics[1][2:], 16)
		fromChainId = a.ChangeID(fromChainId)
		res.FromChainId = (*model.BigInt)(fromChainId)
		res.ToAddress = "0x" + inPair[1].Data[2+24:2+64]
		token, err := a.GetPoolToken(chain, inPair[1].Address)
		if err == nil {
			res.Token = token
		} else {
			log.Error("stargate: cannot get pool token", "chain", chain, "hash", inPair[1].Hash, "pool", inPair[1].Address, "err", err)
			continue
		}
		amount := utils.HexSum(inPair[1].Data[2+64 : 2+2*64])
		convRate, err := a.GetPoolConvertRate(chain, inPair[1].Address)
		if err != nil {
			log.Error("stargate: cannot get pool convert rate", "chain", chain, "hash", inPair[1].Hash, "pool", inPair[1].Address, "err", err)
			continue
		}
		amount.Mul(amount, convRate)
		res.Amount = (*model.BigInt)(amount)

		datas, err := DecodePacketReceivedData(inPair[0].Data)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if len(datas) != 3 {
			log.Error("stargate: invalid PacketReceived log", "chain", chain, "hash", inPair[1].Hash)
			continue
		}
		srcAddress, ok := datas[0].([]byte)
		if !ok {
			log.Error("stargate: invalid PacketReceived log", "chain", chain, "hash", inPair[1].Hash)
			continue
		}
		nonce, ok := datas[1].(uint64)
		if !ok {
			log.Error("stargate: invalid PacketReceived log", "chain", chain, "hash", inPair[1].Hash)
			continue
		}
		d := &InDetail{
			SrcAddress: hexutil.Encode(srcAddress),
			DstAddress: "0x" + inPair[0].Topics[2][26:],
			Nonce:      nonce,
		}
		res.Detail, _ = json.Marshal(d)
		res.MatchTag = strconv.FormatUint(nonce, 10)
		ret = append(ret, res)
	}
	return ret
}

func (a *Stargate) GetPoolToken(chain, pool string) (string, error) {
	p := a.svc.Providers.Get(chain)
	if p == nil {
		return "", fmt.Errorf("providers does not support %v", chain)
	}
	// 0xfc0c546a: token()
	raw, err := p.ContinueCall("", pool, "0xfc0c546a", nil, nil)
	if err != nil {
		return "", err
	}
	return strings.ToLower(common.BytesToAddress(raw).Hex()), nil
}

func (a *Stargate) GetPoolConvertRate(chain, pool string) (*big.Int, error) {
	p := a.svc.Providers.Get(chain)
	if p == nil {
		return nil, fmt.Errorf("providers does not support %v", chain)
	}
	// 0xfeb56b15: convertRate()
	raw, err := p.ContinueCall("", pool, "0xfeb56b15", nil, nil)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(raw), nil
}

func FindParis(events model.Events, sig1, sig2 string) [][2]*model.Event {
	ret := make([][2]*model.Event, 0)
	var i, j int
	for i = 0; i < len(events)-1; i++ {
		if sig1 != events[i].Topics[0] {
			continue
		}
		for j = i + 1; j < len(events); j++ {
			if events[i].Hash != events[j].Hash {
				break
			}
			if events[j].Topics[0] == sig2 {
				ret = append(ret, [2]*model.Event{events[i], events[j]})
				break
			}
		}
		if j >= len(events) {
			return ret
		} else {
			i = j
		}
	}
	return ret
}
