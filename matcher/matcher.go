package matcher

import (
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/lib/pq"
)

const (
	interval  = 40 //扫一次数据库的时间
	batchSize = 10000
)

type Matcher struct {
	svc      *svc.ServiceContext
	projects map[string]model.Matcher
}

func NewMatcher(svc *svc.ServiceContext) *Matcher {
	return &Matcher{
		svc: svc,
		/*projects: map[string]model.Matcher{
			anyswap.NewAnyswapCollector(nil).Name():   NewSimpleInMatcher(svc.ProjectsDao),
			across.NewAcrossCollector().Name():        NewSimpleInMatcher(svc.ProjectsDao),
			celer_bridge.NewCBridgeCollector().Name(): NewSimpleInMatcher(svc.ProjectsDao),
			wormhole.NewWormHoleCollector(nil).Name(): NewSimpleInMatcher(svc.ProjectsDao),
			stargate.NewStargateCollector(nil).Name(): NewSimpleInMatcher(svc.ProjectsDao),
			synapse.NewSynapseCollector(nil).Name():   NewSimpleInMatcher(svc.ProjectsDao),
		},*/
		projects: map[string]model.Matcher{
			"anyswap": NewSimpleInMatcher(svc.ProjectsDao, uint64(7121480)),
			"across":  NewSimpleInMatcher(svc.ProjectsDao, uint64(1547700)),
		},
	}
}

func (m *Matcher) Start() {
	for project, matcher := range m.projects {
		go m.StartMatch(project, matcher)
	}
}

func (m *Matcher) StartMatch(project string, matcher model.Matcher) {
	m.svc.Wg.Add(1)
	defer m.svc.Wg.Done()
	log.Info("matcher start", "project", project)
	timer := time.NewTimer(1 * time.Second)
	cnt := 0
	var LAST = matcher.LastId()
	var last = LAST
	for {
		select {
		case <-m.svc.Ctx.Done():
			log.Info("match svc done", "project", project, "current Id", last)
			return
		case <-timer.C:
			cnt += 1
			latest, err := m.svc.Dao.LatestId(project)
			if err != nil {
				log.Error("get latest id failed", "projet", project, "err", err)
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
				total, matched, err := m.BeginMatch(last+1, right, project, matcher)
				if err != nil {
					log.Error("match job failed", "project", project, "from", last+1, "to", right, "err", err)
				} else {
					last = right
					log.Info("match done", "project", project, "current Id", last, "total", total, "matched crossIn", matched)
				}
			}
			if cnt >= 10 {
				cnt = 0
				last = LAST
			}
		}
		timer.Reset(interval * time.Second)
	}
}

func (m *Matcher) BeginMatch(from, to uint64, project string, matcher model.Matcher) (total int, matched int, err error) {
	var stmt string
	switch matcher.(type) {
	case *SimpleInMatcher:
		//stmt = fmt.Sprintf("select * from %s where project = '%s' and direction = '%s' and id >= $1 and id <= $2 and match_id is null and match_tag not in ('0', '1', '2', '3', '4')", m.svc.Dao.Table(), project, model.InDirection)
		stmt = fmt.Sprintf("select %s from %s where direction = '%s' and id >= $1 and id <= $2 and match_id is null "+
			"", model.ResultRows, project, model.InDirection)
	default:
		panic("invalid matcher")
	}
	var results model.Datas

	if project == "anyswap" {
		stmt_ := fmt.Sprintf("select %s from %s where id >= $1 and id <= $2 and match_id is null ", model.ResultRows, project)
		err = m.svc.Dao.DB().Select(&results, stmt_, from, to)
		if err != nil {
			return
		}

		cnt, errs := matcher.UpdateAnyswapMatchTag(project, results)
		if errs != nil {
			utils.LogError(errs, "./error.log")
		}
		log.Info("update anyswap tag done", "updated", cnt)
	}

	err = m.svc.Dao.DB().Select(&results, stmt, from, to)
	if err != nil {
		return
	}
	shouldUpdates, err := matcher.Match(project, results)
	if err != nil {
		return
	}
	err = m.svc.Dao.Update(project, shouldUpdates)
	if err != nil {
		return
	}
	matched = len(shouldUpdates) / 2
	total = len(results)
	return
}
