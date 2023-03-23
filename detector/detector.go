package detector

import (
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

const (
	interval_fake = 5
	batchSize     = 800
)

type Detector struct {
	svc      *svc.ServiceContext
	projects map[string]model.Detector
}

func NewDetector(svc *svc.ServiceContext, config_path string, startIds map[string]uint64) *Detector {
	return &Detector{
		svc: svc,
		projects: map[string]model.Detector{
			"anyswap": NewSimpleOutDetector(svc, "anyswap", "ethereum", config_path, startIds["anyswap"]),
			"across":  NewSimpleOutDetector(svc, "across", "ethereum", config_path, startIds["across"]),
		},
	}
}

// 检查fake token和chainID
//每2分钟检查一次

func (m *Detector) Start() {
	for project, detector := range m.projects {
		go m.StartDetectOutTx(project, detector)
	}
}

func (m *Detector) StartDetectOutTx(project string, detector model.Detector) {
	m.svc.Wg.Add(1)
	defer m.svc.Wg.Done()
	var last = detector.LastDetectId()
	log.Info("fakeDetector start", "project", project, "start Id", last)
	timer := time.NewTimer(1 * time.Second)
	for {
		t1 := time.Now()
		select {
		case <-m.svc.Ctx.Done():
			log.Info("detect svc done", "project", project, "current Id", last)
			return
		case <-timer.C:
			latest, err := m.svc.Dao.LatestId(project)
			if err != nil {
				log.Error("StartDetectOutTx() get latest id failed", "projet", project, "err", err)
				break
			}
			for last < latest {
				var shouldBreak bool
				select {
				case <-m.svc.Ctx.Done():
					shouldBreak = true
				default:
				}
				if shouldBreak {
					break
				}
				right := utils.Min(latest, last+batchSize)
				err := m.beginDetectOutTx(last+1, right, project, detector)
				if err != nil {
					log.Error("detectFake job failed", "project", project, "from", last+1, "to", right, "err", err)
				} else {
					last = right
					t2 := time.Now()
					log.Info("\ndetectFake done", "project", project, "current Id", last, "batch size", batchSize,
						"time", t2.Sub(t1).String())
				}
			}
		}
		timer.Reset(interval_fake * time.Second)
	}
}

func (m *Detector) beginDetectOutTx(from, to uint64, project string, detector model.Detector) (err error) {
	var stmt string
	switch detector.(type) {
	case *SimpleOutDetector:
		stmt = fmt.Sprintf("select %s from %s where direction = '%s' and id >= $1 and id <= $2", model.ResultRows, project, model.OutDirection)
	default:
		panic("invalid detector")
	}
	var results model.Datas
	err = m.svc.Dao.DB().Select(&results, stmt, from, to)
	if err != nil {
		return
	}
	detector.DetectOutTx(results)
	if err != nil {
		return
	}
	return
}
