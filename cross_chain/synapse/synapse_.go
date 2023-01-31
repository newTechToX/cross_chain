package synapse

/*import (
	"app/model"
	"app/utils"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

var _ model.EventCollector = &Synapse0{}

type Synapse0 struct {
}

func NewSynapseCollector0() *Synapse0 {
	return &Synapse0= {}
}

func (a *Synapse) Name_() string {
	return "Synapse"
}

func (a *Synapse) Contracts_(chain string) []string {
	if _, ok := SynapseContracts[chain]; !ok {
		return nil
	}
	return SynapseContracts[chain]
}

func (a *Synapse) Topics0(chain string) []string {
	return []string{TokenDeposit, TokenDepositAndSwap, TokenMint,
		TokenMintAndSwap, TokenRedeem, TokenRedeemAndRemove, TokenRedeemAndSwap,
		TokenWithdraw, TokenWithdrawAndRemove}
}

func (a *Synapse) Extract_(chain string, events model.Events) model.Results {
	ret := make(model.Results, 0)
	var kappa string
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

		res.ToAddress = "0x" + e.Topics[1][26:]

		if e.Topics[0] == TokenDeposit || e.Topics[0] == TokenDepositAndSwap ||
			e.Topics[0] == TokenRedeem || e.Topics[0] == TokenRedeemAndRemove || e.Topics[0] == TokenRedeemAndSwap {
			if len(e.Topics) < 2 {
				continue
			}

			fromChainId := new(big.Int).Set(utils.GetChainId(chain))
			res.FromChainId = (*model.BigInt)(fromChainId)
			res.Direction = model.OutDirection
			res.Token = "0x" + e.Data[2+64+24:2+128]
			toChainId, _ := new(big.Int).SetString(e.Data[2:2+64], 16)
			res.ToChainId = (*model.BigInt)(toChainId)
			amount, _ := new(big.Int).SetString(e.Data[2+128:2+192], 16)
			res.Amount = (*model.BigInt)(amount)
			var t = crypto.Keccak256Hash([]byte(res.Hash)).String()
			kappa = t
		}
		if e.Topics[0] == TokenMint || e.Topics[0] == TokenMintAndSwap ||
			e.Topics[0] == TokenWithdraw || e.Topics[0] == TokenWithdrawAndRemove {
			res.Direction = model.InDirection
			res.Token = "0x" + e.Data[2+24:2+64]
			amount, _ := new(big.Int).SetString(e.Data[2+64:2+128], 16)
			res.Amount = (*model.BigInt)(amount)
			toChainId := new(big.Int).Set(utils.GetChainId(chain))
			res.ToChainId = (*model.BigInt)(toChainId)
			kappa = e.Topics[len(e.Topics)-1]
		}
		d := &Detail{
			Kappa: kappa,
		}
		detail, err := json.Marshal(d)
		if err == nil {
			res.Detail = detail
		}
		res.MatchTag = kappa
		ret = append(ret, res)
	}
	return ret
}
*/
