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
	"github.com/schollz/progressbar/v3"
	"sync"
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

func (a *Logic) CheckOutTx(project string, datas model.Datas, detected chan int, wg *sync.WaitGroup, limiter chan bool, bar *progressbar.ProgressBar) {
	if datas == nil || len(datas) == 0 {
		return
	}
	defer wg.Done()
	var num = 0
	for _, data := range datas {
		//先检查fakeToken的情况
		if isfake := a.fake_checker.IsFakeToken(project, data); isfake != check_fake.SAFE {
			info := fmt.Sprintf("fake token: fake token in database, chain:%s, address:%s, hash:%s", data.Chain, data.Token, data.Hash)
			utils.SendMail("FAKE TOKEN DETECTED ", info)
			num++
		} else {
			data.IsFakeToken.Scan(1)
			stmt := fmt.Sprintf("update %s set isfaketoken = 0 where id = %d", project, data.Id)
			if _, err := a.svc.Dao.DB().Exec(stmt); err != nil {
				log.Warn("failed to update safe token: ", data.Token, err)
			}
		}

		//查replayOutLogic
		if _, ok := replay.ReplaySupportChains[data.Chain]; !ok {
			continue
		}
		tag, err := a.replayer.ReplayOutTxLogic(project, data)
		if err != nil {
			log.Warn("CheckOutTx(), ", "err", err)
		}
		if tag.TokenProfitError != check_fake.SAFE || tag.FromAddressError != check_fake.SAFE ||
			tag.ToAddressProfit != check_fake.SAFE {
			num++
			info := fmt.Sprintf("%s out tx error: chain:%s, hash:%s, token profit: %d, from profit: %d, to profit: %d",
				project, data.Chain, data.Hash, tag.TokenProfitError, tag.FromAddressError, tag.ToAddressProfit)
			utils.SendMail("OUT TX ERROR DETECTED ", info)
		}

		bar.Add(1)
	}
	detected <- num
	//bar.Add(1)
	<-limiter
	return
}

/*func (a *Logic) CheckOutTx(project string, datas model.Datas, unsafe_tokens_chan chan map[int]model.Datas, wg *sync.WaitGroup, limiter chan bool, bar *progressbar.ProgressBar) {
	if datas == nil || len(datas) == 0 {
		return
	}
	var res_list = make(map[int]model.Datas)
	defer wg.Done()
	for _, data := range datas {
		//先检查fakeToken的情况
		if isfake := a.fake_checker.IsFakeToken(project, data); isfake != check_fake.SAFE {
			res_list[isfake] = append(res_list[isfake], data)
		} else {
			data.IsFakeToken.Scan(1)
			stmt := fmt.Sprintf("update %s set isfaketoken = 0 where id = %data", project, data.Id)
			if _, err := a.svc.Dao.DB().Exec(stmt); err != nil {
				fmt.Println("failed to update safe token: ", data.Token, err)
			}
		}

		//查replayOutLogic
		var tag_chan = make(chan replay.Tags)
		a.replayer.ReplayOutTxLogic(project, data, tag_chan)

		bar.Add(1)
	}
	//bar.Add(1)
	unsafe_tokens_chan <- res_list
	<-limiter
	return
}*/

/*func (a *Logic) CheckOutProfits(project string, datas model.Datas) {
	t1 := time.Now()

	var tag_chan = make(chan replay.Tags)
	for _, data := range datas {
		go a.replayer.ReplayOutTxLogic(project, data, tag_chan)
	}
	t2 := time.Now()
	log.Info("CheckOutTx() done", "time", t2.Sub(t1).String(), "total", len(datas), "currentID", datas[len(datas)-1].Id)
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
}*/
