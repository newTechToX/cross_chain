package replay

import (
	"app/dao"
	"app/model"
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

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

func TestReplayer_CalValue(t *testing.T) {
	r := NewReplayer("../txt_config.yaml")
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
