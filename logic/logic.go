package logic

import (
	"app/logic/aml"
	"app/logic/check_fake"
	"app/logic/check_in"
	"app/logic/replay"
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type Logic struct {
	svc          *svc.ServiceContext
	replayer     *replay.Replayer
	aml          *aml.AML
	fake_checker *check_fake.FakeChecker
	in_checker   *check_in.InChecker
}

func NewLogic(svc *svc.ServiceContext, chain string, config_path string) *Logic {
	c := check_fake.NewChecker(svc, chain, config_path)
	r := replay.NewReplayer(svc, c.Aml(), config_path)
	return &Logic{
		svc:          svc,
		replayer:     r,
		fake_checker: c,
		aml:          c.Aml(),
	}
}

// fake token 和 fake chainId
//chainID的检查还没完成

func (a *Logic) CheckFake(project string, datas model.Datas, unsafe_tokens_chan chan map[int]model.Datas) {
	if datas == nil || len(datas) == 0 {
		return
	}
	var res_list = make(map[int]model.Datas)

	for _, d := range datas {
		if isfake := a.fake_checker.IsFakeToken(project, d); isfake != check_fake.SAFE {
			res_list[isfake] = append(res_list[isfake], d)
		} else {
			stmt := fmt.Sprintf("update %s set isfaketoken = 0", project)
			if _, err := a.svc.Dao.DB().Exec(stmt); err != nil {
				fmt.Println("failed to update safe token: ", d.Token)
			}
		}
	}
	unsafe_tokens_chan <- res_list
	return
}

func (a *Logic) CheckOutTx(project string, datas model.Datas) {
	t1 := time.Now()

	var tag_chan = make(chan replay.Tags)
	for _, data := range datas {
		go a.replayer.ReplayOutTxLogic(project, data, tag_chan)
	}
	t2 := time.Now()
	log.Info("CheckFake() done", "time", t2.Sub(t1).String(), "total", len(datas), "currentID", datas[len(datas)-1].Id)
	return
}

func (a *Logic) CheckDuplicateIn(project string, datas model.Datas) {
	path := "./logic.log"
	t1 := time.Now()
	dup_map := a.in_checker.HasDuplicates(project, datas)
	for k, dup_ids := range dup_map {
		var info string
		for _, id := range dup_ids {
			info += fmt.Sprintf("%d ", id)
		}
		s := fmt.Sprintf("%d has duplicates: %s", k, info)
		utils.LogPrint(s, path)
	}
	t2 := time.Now()
	log.Info("CheckDuplicateIn() done", "time", t2.Sub(t1).String(), "total", len(datas), "currentID", datas[len(datas)-1].Id)
	return
}
