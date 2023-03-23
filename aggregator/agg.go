package aggregator

import (
	crosschain "app/cross_chain"
	"app/model"
	"app/provider"
	"app/svc"
	"app/utils"
	"fmt"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

var (
	BatchSize uint64 = 9000
)

type Aggregator struct {
	svc        *svc.ServiceContext
	chain      string
	provider   *provider.Provider
	collectors []model.Colletcor
}

func NewAggregator(svc *svc.ServiceContext, chain string) *Aggregator {
	p := svc.Providers.Get(chain)
	if p == nil {
		panic(fmt.Sprintf("%v: invalid provider", chain))
	}
	return &Aggregator{
		svc:        svc,
		chain:      chain,
		provider:   p,
		collectors: crosschain.GetCollectors(svc),
	}
}

func (a *Aggregator) Start() {
	for _, c := range a.collectors {
		go a.DoJob(c)
		time.Sleep(time.Second)
	}
}

func (a *Aggregator) DoJob(c model.Colletcor) {
	a.svc.Wg.Add(1)
	defer a.svc.Wg.Done()
	if len(c.Contracts(a.chain)) == 0 {
		return
	}

	timer := time.NewTimer(1 * time.Second)
	last, err := a.getCkpt(c.Name(), c.Contracts(a.chain))
	if err != nil {
		panic(fmt.Sprintf("%v: check failed, %v %v", a.chain, c.Name(), err))
	}
	log.Info("start collector", "chain", a.chain, "project", c.Name(), "last commit", last)
	batchSize := BatchSize
	for {
		select {
		case <-a.svc.Ctx.Done():
			return
		case <-timer.C:
			latest, err := a.provider.LatestNumber()
			if err != nil {
				log.Error("get latest number failed", "chain", a.chain, "err", err.Error())
				break
			}
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
				fetched, err := a.Work(c, last+1, right)
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
		timer.Reset(60 * time.Second)
	}
}

func (a *Aggregator) Work(c model.Colletcor, from, to uint64) (int, error) {
	var totalFetched int
	var results model.Datas

	switch v := c.(type) {
	case model.EventCollector:
		topics0 := v.Topics0(a.chain)
		events, err := a.provider.GetLogs(topics0, from, to)
		if err != nil {
			return 0, err
		}
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
	err := a.svc.Dao.Save(results)
	if err != nil {
		return 0, fmt.Errorf("result save failed: %v", err)
	}
	return totalFetched, nil
}

func (a *Aggregator) getAddressesFirstInvocation(addresses []string) (uint64, error) {
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

func (a *Aggregator) getCkpt(project string, addresses []string) (uint64, error) {
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
