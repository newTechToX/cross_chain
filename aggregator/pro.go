package aggregator

import (
	"app/cross_chain/across"
	"app/cross_chain/anyswap"
	"app/model"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"sort"
)

func (a *Aggregator) StartPro(project string, from, to uint64) {
	var c model.Colletcor
	switch project {
	case "anyswap":
		c = anyswap.NewAnyswapCollector(a.svc)
	case "across":
		c = across.NewAcrossCollector()
	default:
		c = nil
	}
	a.provider = a.svc.Providers.Get(a.chain)
	a.DoJobPro(c, from, to)
}

/*func (a *Aggregator) DoJobPro(c model.Colletcor, from, to uint64) {
	last := from
	latest := to
	//log.Info("start collector_pro", "chain", a.chain, "project", c.Name(), "last commit", last)
	batchSize := BatchSize

	for last < latest {

		right := utils.Min(latest, last+batchSize)
		fetched, err := a.WorkPro(c, last+1, right)

		if err == utils.ErrTooManyRecords {
			batchSize = batchSize / 2
			log.Warn("too many req records", "chain", a.chain, "project", c.Name(), "batch size", batchSize)
		} else if err != nil {
			if err == utils.ErrEtherscanRateLimit {
				log.Warn("etherscan rate limit", "chain", a.chain, "project", c.Name())
			} else {
				log.Error("job failed", "chain", a.chain, "project", c.Name(), "from", last+1, "to", right, "err", err)
			}
		} else {
			last = right
			//log.Info("collect done", "chain", a.chain, "project", c.Name(), "current number", last, "batch size", batchSize)
			if fetched < utils.EtherScanMaxResult*0.8 && batchSize <= 10*utils.EtherScanMaxResult {
				batchSize += 100
			}
		}
	}
}*/

func (a *Aggregator) DoJobPro(c model.Colletcor, from, to uint64) {
	a.svc.Wg.Add(1)
	defer a.svc.Wg.Done()
	/*if len(c.Contracts(a.chain)) == 0 {
		return
	}*/

	last := from
	latest := to
	//log.Info("start collector_pro", "chain", a.chain, "project", c.Name(), "last commit", last)
	batchSize := BatchSize

	for last < latest {
		var shouldBreak bool
		select {
		case <-a.svc.Ctx.Done():
			log.Warn("DoJobPro done")
			shouldBreak = true
		default:
		}
		if shouldBreak {
			break
		}
		right := utils.Min(latest, last+batchSize)
		fetched, err := a.WorkPro(c, last+1, right)

		if err == utils.ErrTooManyRecords {
			batchSize = batchSize / 2
			log.Warn("too many req records", "chain", a.chain, "project", c.Name(), "batch size", batchSize)
		} else if err != nil {
			if err == utils.ErrEtherscanRateLimit {
				log.Warn("etherscan rate limit", "chain", a.chain, "project", c.Name())
			} else {
				log.Error("job failed", "chain", a.chain, "project", c.Name(), "from", last+1, "to", right, "err", err)
			}
		} else {
			last = right
			//log.Info("collect done", "chain", a.chain, "project", c.Name(), "current number", last, "batch size", batchSize)
			if fetched < utils.EtherScanMaxResult*0.8 && batchSize <= 10*utils.EtherScanMaxResult {
				batchSize += 100
			}
		}
	}
}

func (a *Aggregator) WorkPro(c model.Colletcor, from, to uint64) (int, error) {
	var totalFetched int
	var results model.Datas
	addrs := c.Contracts(a.chain)

	if len(addrs) == 0 {
		return 0, nil
	}
	switch v := c.(type) {
	case model.EventCollector:
		topics0 := v.Topics0(a.chain)
		events, err := a.provider.GetLogs(topics0, from, to)
		if err != nil {
			return 0, err
		}
		events = a.filterEvents(c.Name(), events)
		totalFetched = len(events)
		sort.Sort(events)
		results = v.Extract(a.chain, events)
	case model.MsgCollector:
		addrs := c.Contracts(a.chain)
		if len(addrs) == 0 {
			return 0, nil
		}
		selectors := v.Selectors(a.chain)
		calls, err := a.provider.GetCalls(addrs, selectors, from, to)
		if err != nil {
			return 0, err
		}
		totalFetched = len(calls)
		results = v.Extract(a.chain, calls)
	default:
		panic("invalid collector")
	}

	err := a.svc.Dao.Save(results, c.Name())
	if err != nil {
		log.Error("failed to save", "ERROR", err)
		return 0, fmt.Errorf("result save failed: %v", err)
	}
	//println(len(results))
	return totalFetched, nil
}

func (a *Aggregator) filterEvents(project string, events model.Events) model.Events {
	var b model.Events
	for _, event := range events {
		if event.Hash == "0xa86c1701ffece0a07fea6b14a826556099cf0afaf41f4984ef8a343f43a6f0ad" {
			println("0xa86c1701ffece0a07fea6b14a826556099cf0afaf41f4984ef8a343f43a6f0ad")
		}
		if !a.exsit(project, event.Hash, event.Id) {
			b = append(b, event)
			break
		}
	}
	return b
}

func (a *Aggregator) exsit(project, hash string, log_index uint64) bool {
	stmt := fmt.Sprintf("select id from %s where hash = $1 and log_index = $2", project)
	var id uint64
	err := a.svc.Dao.DB().Get(&id, stmt, hash, log_index)
	if err != nil || id == 0 {
		return false
	}
	return true
}
