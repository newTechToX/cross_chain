package detector

import (
	"app/logic"
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"sync"
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

// OutDetector的 Detect 用于检测所有tx的fake token & chainID，将没有match的做二次检测
// 对于fake chainID的检查还没做完

/*func (a *SimpleOutDetector) DetectOutTx(datas model.Datas) int {
	var n, detected = 3, 0
	var size = len(datas) / n
	log.Info("DetectOutTx() begins")
	for i := 0; i < len(datas); i = i + size {
		var right = utils.Min(i+size, len(datas))
		responseChannel := make(chan int, n)

		// 这里在启动goroutine时, 将用来收集结果的局部变量channel也传递进去
		go a.logic.CheckOutTx(a.project, datas[i:right], responseChannel, wg) //, limiter)
	}
	return detected
}*/

func (a *SimpleOutDetector) DetectOutTx(datas model.Datas) int {
	var n, detected = 5, 0
	var size = len(datas) / n
	//var bar = utils.Bar(size, a.project)
	var wg = &sync.WaitGroup{}
	//var limiter = make(chan bool, 10)
	//defer close(limiter)

	responseChannel := make(chan int, n+1)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	/*wgResponse := &sync.WaitGroup{}
	// 启动读取结果的控制器
	go func() {
		// wgResponse计数器+1
		wgResponse.Add(1)
		// 当 responseChannel被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
		wgResponse.Done()
	}()*/

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
	}*/
	return detected
}
