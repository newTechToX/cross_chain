package detector

import (
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	log2 "log"
	"time"
)

const (
	interval_fake = 1 * 60
	batchSize     = 10000
)

type Detector struct {
	svc      *svc.ServiceContext
	projects map[string]model.Detector
}

func NewDetector(svc *svc.ServiceContext, config_path string) *Detector {
	return &Detector{
		svc: svc,
		projects: map[string]model.Detector{
			"anyswap": NewSimpleOutDetector(svc, "ethereum", config_path, uint64(6972699)),
			"across":  NewSimpleOutDetector(svc, "ethereum", config_path, uint64(1)),
		},
	}
}

// 检查fake token和chainID
//每2分钟检查一次

func (m *Detector) Start() {
	for project, detector := range m.projects {
		go m.StartDetectFake(project, detector)
	}
}

func (m *Detector) StartDetectFake(project string, detector model.Detector) {
	m.svc.Wg.Add(1)
	defer m.svc.Wg.Done()
	var last = detector.LastId()
	log2.SetPrefix("StartDetectFake()")
	log.Info("fakeDetector start", "project", project, "start Id", last)
	timer := time.NewTimer(1 * time.Second)
	for {
		t1 := time.Now()
		select {
		case <-m.svc.Ctx.Done():
			log2.SetPrefix("StartDetectFake()")
			log.Info("detect svc done", "project", project, "current Id", last)
			return
		case <-timer.C:
			latest, err := m.svc.Dao.LatestId(project)
			if err != nil {
				log.Error("StartDetectFake() get latest id failed", "projet", project, "err", err)
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
				fake, err := m.beginDetectFake(last+1, right, project, detector)
				if err != nil {
					log2.SetPrefix("StartDetectFake()")
					log.Error("detectFake job failed", "project", project, "from", last+1, "to", right, "err", err)
				} else {
					last = right
					t2 := time.Now()
					log2.SetPrefix("StartDetectFake()")
					log.Info("detectFake done", "project", project, "current Id", last, "batch size", batchSize,
						"fake", fake, "time", t2.Sub(t1).String())
				}
			}
		}
		timer.Reset(interval_fake * time.Second)
	}
}

func (m *Detector) beginDetectFake(from, to uint64, project string, detector model.Detector) (fake int, err error) {
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
	fake = detector.DetectFake(project, results)
	if err != nil {
		return
	}
	return
}
