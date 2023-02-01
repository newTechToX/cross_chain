package logic

import (
	"app/dao"
	"app/logic/replay"
	"app/model"
	"fmt"
	"testing"
)

func TestLogic_ReplayOutTxLogic(t *testing.T) {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")

	a := &Logic{}
	re := &replay.Replayer{}
	a.replayer = re.NewReplayer()
	//hash := "0x4f2eb92a2a9a21bd0c19eab7b4dd3ff4cea4979b70ea4cf56fe20a6e14f73bbd"
	id := 4167038
	//stmt := fmt.Sprintf("select * from anyswap where hash = '%s'", hash)
	stmt := fmt.Sprintf("select * from anyswap where id = %d", id)

	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}

	err = a.replayOutTxLogic("anyswap", datas)
}

func Test1(t *testing.T) {
	u := map[int]int{
		1: 1,
		2: 2,
	}
	x := 3
	if _, ok := u[4]; ok || x == 2 {
		println("ok")
	}
}

func TestIf(t *testing.T) {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	a := &Logic{}
	re := &replay.Replayer{}
	a.replayer = re.NewReplayer()
	stmt := fmt.Sprintf("select * from anyswap where direction='out' and (chain='ethereum' or chain='bsc') and isfaketoken is null limit 12000")
	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}
	println(len(datas))

	i, size := 0, 600
	for i = 0; i+size < len(datas); i = i + size {
		go ifFromTransferOnly1(datas[i : i+size])
	}
	ifFromTransferOnly1(datas[i:])
}

func ifTokenSourceOnly1(datas []*model.Data) {
	re := &replay.Replayer{}
	r := re.NewReplayer()
	for i, d := range datas {
		tx, err := r.Replay(d)
		if err != nil {
			fmt.Println(err)
		}
		for _, b := range tx.BalanceChanges {
			if b.Account == d.Token && len(b.Assets) > 1 {
				fmt.Println(d.Hash)
			}
		}
		if i%50 == 0 && i != 0 {
			println("done: ", i)
		}
	}
	println("all done")
}

func ifFromTransferOnly1(datas []*model.Data) {
	re := &replay.Replayer{}
	r := re.NewReplayer()
	for i, d := range datas {
		tx, err := r.Replay(d)
		if err != nil {
			fmt.Println(err)
		}
		for _, b := range tx.BalanceChanges {
			if b.Account == d.FromAddress && len(b.Assets) > 1 {
				fmt.Println(d.Hash)
			}
		}
		if i%50 == 0 && i != 0 {
			println("done: ", i)
		}
	}
	println("all done")
}

func selectData(id int) []*model.Data {
	d := dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	//stmt := fmt.Sprintf("select * from anyswap where direction='out' and (chain='ethereum' or chain='bsc') and isfaketoken is null limit 12000")
	stmt := fmt.Sprintf("select * from anyswap where id = %d", id)
	var datas = []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}
	return datas
}

func TestLogic_getPreviousToken(t *testing.T) {
	id := 4266976
	a := &Logic{}
	re := &replay.Replayer{}
	a.replayer = re.NewReplayer()
	datas := selectData(id)
	tx, err := a.replayer.Replay(datas[0])
	if err != nil {
		fmt.Println(err)
	}

	token := "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	res := a.getPreviousToken(token, tx.BalanceChanges)
	for key := range res {
		println(key)
	}
}

func TestCheckFromWithSwap(t *testing.T) {
	id := 4266362
	a := &Logic{}
	re := &replay.Replayer{}
	a.replayer = re.NewReplayer()
	datas := selectData(id)
	data := datas[0]
	tx, err := a.replayer.Replay(datas[0])
	if err != nil {
		fmt.Println(err)
	}
	ETH := "0x2170ed0880ac9a755fd29b2688956bd959f933f8"
	p := a.getPreviousToken(ETH, tx.BalanceChanges)

	for _, e := range tx.BalanceChanges {
		if e.Account == data.ToAddress {
			for _, ee := range e.Assets {
				flag, tokens := a.checkFrom_Token(p, ee.Address, data.Amount.String(), e)
				println(flag)
				println(len(tokens))
			}
		}
	}
}
