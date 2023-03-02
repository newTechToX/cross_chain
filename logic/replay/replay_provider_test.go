package replay

import (
	"app/config"
	"app/dao"
	"app/logic/aml"
	"app/model"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"os"
	"strconv"
	"testing"
)

var cfg config.Config
var srvCtx *svc.ServiceContext

func init() {
	config.LoadCfg(&cfg, "../../config.yaml")
	srvCtx = svc.NewServiceContext(context.Background(), &cfg)
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	chainbase.SetupLimit(10)
}

/*func TestReplay_Replay_tx(t *testing.T) {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	re := &Replayer{}
	r := re.NewReplayer()
	hash := "0x46749dfabb21ab7199c96ea5284ad59346a62676962af9a69f7c60ada6524c36"
	stmt := fmt.Sprintf("select * from anyswap where hash = '%s'", hash)
	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}

	tx, err := r.Replay(datas[0])
	if err != nil {
		println(err)
	}

}*/

func TestLogic_ReplayOutTxLogic(t *testing.T) {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	a := NewReplayer(nil, nil, "../txt_config.yaml")
	//hash := "0x4f2eb92a2a9a21bd0c19eab7b4dd3ff4cea4979b70ea4cf56fe20a6e14f73bbd"
	id := 7585624
	//stmt := fmt.Sprintf("select * from anyswap where hash = '%s'", hash)
	stmt := fmt.Sprintf("select %s from anyswap where id = %d", model.ResultRows, id)

	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}

	a.svc = srvCtx
	a.aml = aml.NewAML("../txt_config.yaml")
	tag, err := a.ReplayOutTxLogic("anyswap", datas[0])
	fmt.Println(tag)
}

func TestReplayer_CalValue(t *testing.T) {
	r := NewReplayer(nil, nil, "../txt_config.yaml")
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	hash := "0xc33f6c406f1172c01d0b987237624f2cbe1021fe721da0d2fb07b31553edb684"
	stmt := fmt.Sprintf("select * from anyswap where hash = '%s'", hash)
	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}
	tx, err := r.Replay(datas[0])
	if err != nil {
		println(err)
	}
	weth := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	for _, e := range tx.BalanceChanges {
		if e.Account == weth {
			ba := r.CalAmount(e)
			fmt.Println(ba)
		}
	}
}

func Test_Cal_bigFloat(t *testing.T) {
	a := "-107097304672796987"
	aa, _ := new(big.Float).SetPrec(uint(256)).SetString(a)
	d := "18"
	dec, _ := strconv.ParseFloat("1e"+d, 64)
	denominator := big.NewFloat(dec)
	denominator1 := aa.Quo(aa, denominator)
	ta := "-563427605274940375"
	tt, _ := new(big.Float).SetPrec(uint(256)).SetString(ta)
	denominator2 := tt.Quo(tt, denominator)

	denominator1.Add(denominator1, denominator2)
	fmt.Println(denominator1)
}

func TestLogic_getPreviousToken(t *testing.T) {
	id := 4266976
	a := NewReplayer(nil, nil, "./txt_config.yaml")
	datas := selectData("", id)
	tx, err := a.Replay(datas[0])
	if err != nil {
		fmt.Println(err)
	}

	token := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	res := a.getPreviousToken(token, tx.BalanceChanges)
	for key := range res {
		println(key)
	}
}

func selectData(st string, id int) []*model.Data {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	//stmt := fmt.Sprintf("select * from anyswap where direction='out' and (chain='ethereum' or chain='bsc') and isfaketoken is null limit 12000")
	var stmt = st
	if stmt == "" {
		stmt = fmt.Sprintf("select * from anyswap where id = %d", id)
	}
	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}
	println(len(datas))
	return datas
}

func TestCheckFromWithSwap(t *testing.T) {
	id := 4266362
	a := NewReplayer(nil, nil, "./txt_config.yaml")
	datas := selectData("", id)
	data := datas[0]
	tx, err := a.Replay(datas[0])
	if err != nil {
		fmt.Println(err)
	}
	ETH := "0x2170ed0880ac9a755fd29b2688956bd959f933f8"
	p := a.getPreviousToken(ETH, tx.BalanceChanges)

	for _, e := range tx.BalanceChanges {
		if e.Account == data.ToAddress {
			for _, ee := range e.Assets {
				flag, tokens := a.checkFromToken(p, ee.Address, data.Amount.String(), e)
				println(flag)
				println(len(tokens))
			}
		}
	}
}

func TestLogic_TokenProfitError1(t *testing.T) {
	id := 338939

	stmt := "select * from anyswap where token_profit_error is not null"
	datas := selectData(stmt, id)

	i, size := 0, 2500
	for i = 0; i < len(datas)-2*size; i = i + size {
		go test_token_error_1(datas[i : i+size])
	}
	println(i)
	test_token_error_1(datas[i:])

}

func test_token_error_1(datas []*model.Data) {
	a := NewReplayer(nil, nil, "./txt_config.yaml")

	for _, d := range datas {
		tx, err := a.Replay(d)
		if err != nil {
			fmt.Println(err)
			continue
		}

		real_token := a.getRealToken("", "", d.Token, tx.BalanceChanges)

		for _, value := range real_token {
			x := new(model.BigInt).SetString(value.Amount, 10)
			if x.Cmp(d.Amount) >= 0 {
				continue
			} else {
				stmt := fmt.Sprintf("update anyswap set token_profit_error = 2 where id = %d", d.Id)
				da := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
				if _, err := da.DB().Exec(stmt); err != nil {
					fmt.Println(err)
				} else {
					println(d.Hash)
				}
			}
		}
	}
	println("all done")
}
