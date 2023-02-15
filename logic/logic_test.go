package logic

import (
	"app/dao"
	"app/logic/replay"
	"app/model"
	"fmt"
	"testing"
)

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
	a.replayer = replay.NewReplayer(nil, nil, "./txt_config.yaml")
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
	r := replay.NewReplayer(nil, nil, "./txt_config.yaml")
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
	r := replay.NewReplayer(nil, nil, "./txt_config.yaml")
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
