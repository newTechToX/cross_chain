package aggregator

import (
	"app/cross_chain/across"
	"app/cross_chain/anyswap"
	"app/model"
	"app/svc"
	"app/utils"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"sort"
)

func (a *Aggregator) StartPro(svc *svc.ServiceContext, project string, from, to uint64) {
	var c model.Colletcor
	switch project {
	case "anyswap":
		c = anyswap.NewAnyswapCollector(svc)
	case "across":
		c = across.NewAcrossCollector()
	default:
		c = nil
	}
	a.DoJobPro(c, from, to)
}

func (a *Aggregator) DoJobPro(c model.Colletcor, from, to uint64) {
	a.svc.Wg.Add(1)
	defer a.svc.Wg.Done()
	if len(c.Contracts(a.chain)) == 0 {
		return
	}

	last := from
	latest := to
	log.Info("start collector_pro", "chain", a.chain, "project", c.Name(), "last commit", last)
	batchSize := BatchSize
	for last < latest {
		var shouldBreak bool
		select {
		case <-a.svc.Ctx.Done():
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
			log.Info("collect done", "chain", a.chain, "project", c.Name(), "current number", last, "batch size", batchSize)
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
		events, err := a.provider.GetLogs(addrs, topics0, from, to)
		if err != nil {
			return 0, err
		}
		totalFetched = len(events)
		sort.Sort(events)
		results = v.Extract(a.chain, events)
	case model.MsgCollector:
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
	err := a.svc.Dao.Save(results)
	if err != nil {
		return 0, fmt.Errorf("result save failed: %v", err)
	}
	return totalFetched, nil
}

func (a *Aggregator) getAddressesFirstInvocationPro(addresses []string) (uint64, error) {
	nums := make([]uint64, 0)
	for _, addr := range addresses {
		n, err := a.provider.GetContractFirstInvocation(addr)
		if err != nil {
			log.Error("get address first invoke failed", "chain", a.chain, "address", addr, "err", err.Error())
		}
		if n != 0 {
			nums = append(nums, n)
		}
	}
	if len(nums) == 0 {
		return 0, nil
	}
	return utils.Min(nums...), nil
}

func (a *Aggregator) getCkptPro(project string, addresses []string) (uint64, error) {
	last, err := a.svc.Dao.LastUpdate(a.chain, project)
	if err != nil {
		return 0, err
	}
	if last == 0 {
		last, err = a.getAddressesFirstInvocation(addresses)
		if err != nil {
			return 0, err
		}
	}
	if last == 0 {
		log.Error("addr first call is 0", "chain", a.chain, "project", project)
		last, err = a.provider.LatestNumber()
		if err != nil {
			return 0, err
		} else {
			last = last - 1000000
		}
	}
	// latest, err := a.provider.LatestNumber()
	// if err != nil {
	// 	return 0, err
	// }
	// last = utils.Max(last, latest-1000000)
	return last, nil
}
