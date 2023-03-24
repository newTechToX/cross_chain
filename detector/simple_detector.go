package detector

import (
	"app/logic"
	"app/logic/check_fake"
	"app/logic/replay"
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	log2 "log"
)

type SimpleOutDetector struct {
	svc      *svc.ServiceContext
	logic    *logic.Logic
	start_id uint64
	project  string
}

var _ model.Detector = &SimpleOutDetector{}

func NewSimpleOutDetector(svc *svc.ServiceContext, project, chain, config_path string, start_id uint64) *SimpleOutDetector {
	return &SimpleOutDetector{
		svc:   svc,
		logic: logic.NewLogic(svc, chain, config_path),
		//matcher: matcher.NewMatcher(svc),
		start_id: start_id,
		project:  project,
	}
}

func (a *SimpleOutDetector) LastDetectId() uint64 {
	stmt := fmt.Sprintf("select max(id) from %s where direction = 'out' and id>%d and isfaketoken = 0", a.project, a.start_id)
	type ID struct {
		Id uint64 `db:"max"`
	}
	var id = ID{a.start_id}
	if err := a.svc.Dao.DB().Get(&id, stmt); err != nil {
		log.Warn("failed to get undetected id", "project", a.project, "ERROR", err)
	} else {
		a.start_id = id.Id
	}
	return a.start_id
}

func (a *SimpleOutDetector) DetectOutTx(datas model.Datas) {
	if datas == nil || len(datas) == 0 {
		return
	}

	log.Info("simpleOutDetector.DetectOutTx() begins")
	var detected = 0
	for _, data := range datas {
		//t1 := time.Now()

		//先检查fakeToken的情况
		if isfake := a.logic.IsFakeToken(a.project, data); isfake != check_fake.SAFE {
			info := fmt.Sprintf("fake token: fake token in database, chain:%s, address:%s, hash:%s", data.Chain, data.Token, data.Hash)
			utils.SendMail(fmt.Sprintf("%s FAKE TOKEN DETECTED", data.Project), info)
			detected++
		} else {
			data.IsFakeToken.Scan(1)
			stmt := fmt.Sprintf("update %s set isfaketoken = 0 where id = %d", a.project, data.Id)
			if _, err := a.svc.Dao.DB().Exec(stmt); err != nil {
				log.Warn("failed to update safe token: ", data.Token, err)
			}
		}

		//查replayOutLogic
		if _, ok := replay.ReplaySupportChains[data.Chain]; !ok {
			continue
		}
		tag, err := a.logic.ReplayOutTxLogic(a.project, data)
		if err != nil {
			log2.SetPrefix("CheckOutTx()")
			utils.LogPrint(err.Error(), "../logic.log")
		}
		if tag.TokenProfitError != replay.SAFE || tag.ToAddressProfit == replay.TOKEN_PROFIT_MINUS_AMOUNT {
			detected++
			info := fmt.Sprintf("%s out tx error: chain:%s, hash:%s, token profit: %d, from profit: %d, to profit: %d",
				a.project, data.Chain, data.Hash, tag.TokenProfitError, tag.FromAddressError, tag.ToAddressProfit)
			utils.SendMail("OUT TX ERROR DETECTED ", info)
		}
	}
	log.Info("simpleOutDetector.DetectOutTx() detected ", detected)
	return
}

// OutDetector的 Detect 用于检测所有tx的fake token & chainID，将没有match的做二次检测
// 对于fake chainID的检查还没做完

/*func (a *SimpleOutDetector) DetectOutTx(datas model.Datas) int {
	var n, detected = 3, 0
	var wg = &sync.WaitGroup{}
	var size = len(datas) / n
	if size == 0 {
		size = len(datas)
		n = 1
	}
	log.Info("DetectOutTx() begins")
	for i := 0; i < len(datas); i = i + size {
		var right = utils.Min(i+size, len(datas))
		responseChannel := make(chan int, size)

		// 这里在启动goroutine时, 将用来收集结果的局部变量channel也传递进去
		go a.detectOutTx(datas[i:right], responseChannel, wg) //, limiter)
	}
	return detected
}*/

/*func (a *SimpleOutDetector) DetectOutTx(datas model.Datas) int {
	var n, detected = 5, 0
	var size = len(datas) / n
	//var bar = utils.Bar(size, a.project)
	var wg = &sync.WaitGroup{}
	//var limiter = make(chan bool, 10)
	//defer close(limiter)

	responseChannel := make(chan int, n+1)

	log.Info("DetectOutTx() begins")
	for i := 0; i < len(datas); i = i + size {
		var right = utils.Min(i+size, len(datas))
		// 计数器+1
		wg.Add(1)
		//limiter <- true
		// 这里在启动goroutine时, 将用来收集结果的局部变量channel也传递进去
		go a.logic.CheckOutTx(a.project, datas[i:right], responseChannel, wg) //, limiter)
		detected += <-responseChannel
	}

	// 等待所以协程执行完毕
	wg.Wait() // 当计数器为0时, 不再阻塞
	// 关闭接收结果channel
	close(responseChannel)
	// 等待wgResponse的计数器归零
	//wgResponse.Wait()

	/*err := bar.Close()
	if err != nil {
		log.Warn("DetectOutTx(): Failed to close bar", "Error", err)
	}
	return detected
}*/
