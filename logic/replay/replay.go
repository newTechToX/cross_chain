package replay

import (
	"app/model"
	"app/utils"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Replayer struct {
	url string
}

func (a *Replayer) NewReplayer() *Replayer {
	return &Replayer{
		url: "https://api.blocksec.com/v1/phalcon/simulate/hash",
	}
}

func (a *Replayer) Replay(data *model.Data) (*SimulatedTxn, error) {
	var tx = &SimulatedTxn{}

	type Body struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	body := Body{}

	time.Sleep(500 * time.Millisecond)
	buffer, err := a.replay(data)
	er := json.Unmarshal(buffer, &body)
	if er != nil {
		return tx, er
	}
	if body.Code == 1 {
		err = fmt.Errorf("1")
		return tx, err
	}
	if err != nil {
		return tx, err
	}
	var simulated_tx = &RawMsg{}
	err = json.Unmarshal(buffer, simulated_tx)
	if len(simulated_tx.Data.Txns) == 0 {
		return nil, err
	}
	tx = simulated_tx.Data.Txns[0]
	if err != nil {
		return tx, err
	}
	if len(tx.BalanceChanges) == 0 {
		println("no balance changes ", data.Hash, data.Chain)
	}

	return tx, err
}

func (a *Replayer) replay(data *model.Data) ([]byte, error) {
	var err error
	token, err := os.ReadFile("./logic/replay/token.txt")
	if err != nil {
		token, err = os.ReadFile("./replay/token.txt")
		if err != nil {
			token, err = os.ReadFile("./token.txt")
			if err != nil {
				log.Error("failed to open token.txt")
				return nil, err
			}
		}
	}

	hash := []string{data.Hash}
	id, _ := strconv.Atoi(utils.GetChainId(data.Chain).String())
	b, _ := json.Marshal(map[string]interface{}{
		"chainID":  id,
		"block":    data.Number,
		"position": data.Index,
		"bundle":   hash,
	})
	req, err := http.NewRequest("POST", a.url, strings.NewReader(string(b)))
	req.Header.Set("Access-Token", string(token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Error("failed request ", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (a *Replayer) CalAmount(balance *SimAccountBalance) *big.Float {
	var ba, _ = new(big.Float).SetPrec(uint(256)).SetString("0")
	for _, asset := range balance.Assets {
		var value = a.GetAmount(asset)
		ba.Add(ba, value)
	}
	return ba
}

func (a *Replayer) GetAmount(asset *SimAsset) *big.Float {
	ret, _ := new(big.Float).SetPrec(uint(256)).SetString(asset.Amount)
	dec := fmt.Sprintf("1e%d", asset.Decimals)
	denominator, _ := new(big.Float).SetString(dec)
	res := ret.Quo(ret, denominator)
	return res
}
