package replay

import (
	"app/config"
	"app/logic/aml"
	"app/model"
	"app/svc"
	"app/utils"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type Replayer struct {
	svc *svc.ServiceContext
	aml *aml.AML
	url string
	key string
}

func NewReplayer(svc *svc.ServiceContext, aml *aml.AML, config_path string) *Replayer {
	var cfg Config
	if config_path == "" {
		config_path = "../txt_config.yaml"
	}
	config.LoadCfg(&cfg, config_path)
	return &Replayer{
		svc: svc,
		aml: aml,
		url: "https://api.blocksec.com/v1/phalcon/simulate/hash",
		key: cfg.ReplayAccessToken,
	}
}

func (a *Replayer) Replay(data *model.Data) (*SimulatedTxn, error) {
	var tx = &SimulatedTxn{}

	type Body struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	body := Body{}

	//time.Sleep(6 * time.Second)
	buffer, err := a.replay(data)
	if err != nil {
		return nil, err
	}
	er := json.Unmarshal(buffer, &body)
	if er != nil || body.Msg != "OK" {
		er = fmt.Errorf("msg: %s, hash: %s, error: %s", body.Msg, data.Hash, er)
		return nil, er
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
		println("Replay(), no balance changes ", data.Hash, data.Chain)
	}

	return tx, err
}

func (a *Replayer) replay(data *model.Data) ([]byte, error) {
	var err error
	hash := []string{data.Hash}
	id, _ := strconv.Atoi(utils.GetChainId(data.Chain).String())
	b, _ := json.Marshal(map[string]interface{}{
		"chainID":  id,
		"block":    data.Number,
		"position": data.TxIndex,
		"bundle":   hash,
	})
	req, err := http.NewRequest("POST", a.url, strings.NewReader(string(b)))
	req.Header.Set("Access-Token", a.key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Error("replay(), failed request ", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (a *Replayer) CalAmount(balance *SimAccountBalance) *model.BigFloat {
	var ba, _ = new(big.Float).SetPrec(uint(256)).SetString("0")
	for _, asset := range balance.Assets {
		var value = a.GetFloatAmount(asset.Amount, asset.Decimals)
		ba.Add(ba, (*big.Float)(value))
	}
	return (*model.BigFloat)(ba)
}

func (a *Replayer) GetFloatAmount(Amount string, Decimals int) *model.BigFloat {
	ret, _ := new(big.Float).SetPrec(uint(256)).SetString(Amount)
	dec := fmt.Sprintf("1e%d", Decimals)
	denominator, _ := new(big.Float).SetString(dec)
	res := ret.Quo(ret, denominator)
	return (*model.BigFloat)(res)
}
