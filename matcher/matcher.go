package matcher

import (
	"app/aggregator"
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/lib/pq"
)

const (
	interval  = 40 //扫一次数据库的时间
	batchSize = 10000
)

var Projects = []string{
	"anyswap", "across", //"synapse",
}

type Matcher struct {
	svc      *svc.ServiceContext
	projects map[string]model.Matcher
}

func NewMatcher(svc *svc.ServiceContext, startIds map[string]uint64) *Matcher {
	var projects = make(map[string]model.Matcher)
	for _, project := range Projects {
		projects[project] = NewSimpleInMatcher(project, svc.ProjectsDao, startIds[project])
	}
	return &Matcher{
		svc:      svc,
		projects: projects,
	}
}

func (m *Matcher) Start() {
	for _, matcher := range m.projects {
		go m.StartMatch(matcher)
	}
}

func (m *Matcher) StartMatch(matcher model.Matcher) {
	m.svc.Wg.Add(1)
	defer m.svc.Wg.Done()
	timer := time.NewTimer(1 * time.Second)
	//var last = matcher.LastUnmatchId()
	var last = matcher.LastUnmatchId()
	log.Info("matcher start", "project", matcher.Project(), "Start ID", last)
	for {
		select {
		case <-m.svc.Ctx.Done():
			log.Info("match svc done", "project", matcher.Project(), "current Id", last)
			return
		case <-timer.C:
			latest, err := m.svc.Dao.LatestId(matcher.Project())
			if err != nil {
				log.Error("get latest id failed", "projet", matcher.Project(), "err", err)
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
				if last > 1000 {
					last -= 1000
				} else {
					last = 1
				}
				total, matched, err := m.BeginMatch(last+1, right, matcher.Project(), matcher)
				if err != nil {
					log.Error("match job failed", "project", matcher.Project(), "from", last+1, "to", right, "err", err)
				} else {
					last = right
					log.Info("match done", "project", matcher.Project(), "current Id", last, "total", total, "matched crossIn", matched)
				}
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
		stmt = fmt.Sprintf("select %s from %s where direction = '%s' and id >= $1 and id <= $2 and match_id is null order by id asc",
			model.ResultRows, project, model.InDirection)
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

		cnt, errs := matcher.UpdateAnyswapMatchTag(results)
		if errs != nil {
			utils.LogError(errs, "./error.log")
		}
		log.Info("update anyswap tag done", "updated", cnt)
	}

	err = m.svc.Dao.DB().Select(&results, stmt, from, to)
	if err != nil {
		return
	}
	shouldUpdates_1, unmatches_map, err := matcher.Match(results)
	if err != nil {
		return
	}
	err = m.svc.Dao.Update(project, shouldUpdates_1)
	if err != nil {
		return
	}
	shouldUpdates_2 := m.ProcessUnmatch(project, unmatches_map, matcher)
	err = m.svc.Dao.Update(project, shouldUpdates_2)
	if err != nil {
		return
	}
	matched = (len(shouldUpdates_1) + len(shouldUpdates_2)) / 2
	total = len(results)
	return
}

type Block struct {
	MAX   uint64 `db:"max"`
	MIN   uint64 `db:"min"`
	Chain string `db:"chain"`
}

type Blocks []*Block

func (m *Matcher) ProcessUnmatch(project string, unmatches_map map[string]model.Datas, matcher model.Matcher) (shouldUpdates model.Datas) {
	blocks := m.extractUnmatchInfo(project, unmatches_map)

	var wg sync.WaitGroup
	for _, b := range blocks {
		wg.Add(1)
		go func(svc *svc.ServiceContext, project, chain string, from, to uint64) {
			defer wg.Done()
			log.Info("matcher.ProcessUnmatch start ", "Chain", chain, "From", from, "To", to)
			agg := aggregator.NewAggregator(svc, chain)
			agg.StartPro(project, from, to)
		}(m.svc, project, b.Chain, b.MIN, b.MAX)

		/*cmd := exec.Command("./pro", "-name", "anyswap", "-from", from_block, "-to", to_block, "-chain", b.Chain)
		dd, _ := cmd.Output()
		log.Info("Process Unmatch done")
		if string(dd) != "" {
			fetched, _ := strconv.Atoi(string(dd))
			return fetched
		}*/
	}
	wg.Wait()

	var pre_unmatches model.Datas
	for _, datas := range unmatches_map {
		pre_unmatches = append(pre_unmatches, datas...)
	}

	shouldUpdates, still_unmatches, err := matcher.Match(pre_unmatches)
	if err != nil {
		return
	}
	for _, un := range still_unmatches {
		matcher.SendMail("UNMATCH", un)
	}
	log.Info("matcher.ProcessUnmatch done", "total ", len(pre_unmatches), "matched ", len(shouldUpdates), "unmatch ", len(still_unmatches))
	return
}

func (a *Matcher) extractUnmatchInfo(project string, unmatches_map map[string]model.Datas) (unmatch_blocks Blocks) {
	for chainId, datas := range unmatches_map {
		s := fmt.Sprintf("select max(block_number), min(block_number), chain from %s where id >= $1 and direction = 'out' and from_chain = %s group by chain", project, chainId)
		var blocks Blocks
		err := a.svc.Dao.DB().Select(&blocks, s, datas[0].Id-1000)
		if err != nil {
			log.Warn("Matcher.extractUnmatchInfo failed ", "Error", err)
			continue
		}
		unmatch_blocks = append(unmatch_blocks, blocks...)
	}
	return
}
