package logic

import (
	"app/config"
	"app/dao"
	"app/logic/check_in"
	"app/logic/replay"
	"app/model"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"os"
	"testing"
)

var cfg config.Config
var srvCtx *svc.ServiceContext

func init() {
	config.LoadCfg(&cfg, "../config.yaml")
	srvCtx = svc.NewServiceContext(context.Background(), &cfg)
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	chainbase.SetupLimit(10)
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

var d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")

func TestIf(t *testing.T) {
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

func TestLogic_CheckDuplicateIn(t *testing.T) {
	stmt := fmt.Sprintf("select %s from across where direction = 'in' and match_id is null", model.ResultRows)
	var datas model.Datas
	d.DB().Select(&datas, stmt)
	a := NewLogic(srvCtx, "bsc", "./txt_config.yaml")
	a.in_checker = check_in.NewInChecker(d)

	i := 0
	size := len(datas) / 50
	for ; i < len(datas)-2*size; i = i + size {
		go a.CheckDuplicateIn("across", datas[i:i+size])
	}
	a.CheckDuplicateIn("across", datas[i:])
}
