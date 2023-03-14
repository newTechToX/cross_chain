package etherscan

import (
	"app/model"
	"app/utils"
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/time/rate"
)

const (
	MaxBlockNumber  = 999999999
	DefaultPageSize = 10000

	normalTxApi      = "%s/api?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s&apikey=%s"
	internalTxApi    = "%s/api?module=account&action=txlistinternal&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s&apikey=%s"
	logWithTopicsApi = "%s/api?module=logs&action=getLogs&fromBlock=%d&toBlock=%d&topic0=%s&page=%d&offset=%d&apikey=%s"
	logApi           = "%s/api?module=logs&action=getLogs&fromBlock=%d&toBlock=%d&page=%d&offset=%d&apikey=%s"

	latestNumApi = "%s/api?module=proxy&action=eth_blockNumber&apiKey=%s"

	noTransactionsFound = "No transactions found"
	noRecordsFound      = "No records found"
)

type Option struct {
	Page       int
	PageSize   int
	Asc        bool
	StartBlock int64
	EndBlock   int64
}

type EtherscanProvider struct {
	baseUrl string
	proxy   string
	// key pool
	apiKeys []string
	keyIter uint
	l       sync.Mutex
	limiter *rate.Limiter
}

func NewEtherScanProvider(baseUrl string, apiKeys []string, proxy string, rateLimit int) *EtherscanProvider {
	return &EtherscanProvider{
		baseUrl: strings.TrimRight(baseUrl, "/"),
		proxy:   proxy,
		apiKeys: apiKeys,
		limiter: rate.NewLimiter(rate.Limit(rateLimit), 1),
	}
}

func (p *EtherscanProvider) LatestNumber() (uint64, error) {
	p.limiter.Wait(context.Background())
	url := fmt.Sprintf(latestNumApi, p.baseUrl, p.nextKey())
	log.Debug("invoke etherscan", "url", url)
	var resp GethResponse[string]
	if err := utils.HttpGetObjectWithProxy(url, p.proxy, &resp); err != nil {
		return 0, fmt.Errorf("%v, url: %v", err, url)

	}
	return utils.ParseStrToUint64(resp.Result), nil
}

func (p *EtherscanProvider) GetLogs(topics0 []string, from, to uint64) (model.Events, error) {
	ret := make(model.Events, 0)
	for _, topic0 := range topics0 {
		page := 1
		for {
			if page*DefaultPageSize > utils.EtherScanMaxResult {
				return nil, utils.ErrTooManyRecords
			}
			rawLogs, err := p.getLogs(topic0, Option{
				Page:       page,
				PageSize:   DefaultPageSize,
				StartBlock: int64(from),
				EndBlock:   int64(to),
			})
			if err != nil {
				return nil, err
			}
			for _, l := range rawLogs {
				log := model.Event{
					Ts:      time.Unix(int64(utils.ParseStrToUint64(l.Timestamp)), 0),
					Number:  utils.ParseStrToUint64(l.BlockNumber),
					Index:   utils.ParseStrToUint64(l.Index),
					Hash:    l.Hash,
					Id:      utils.ParseStrToUint64(l.LogIndex),
					Address: l.Address,
					Topics:  l.Topics,
					Data:    l.Data,
				}
				ret = append(ret, &log)
				if log.Number == uint64(16786623) {
					println(log.Hash)
				}
			}
			if len(rawLogs) < DefaultPageSize {
				break
			}
			page++
		}
	}

	return ret, nil
}

func (p *EtherscanProvider) GetCalls(addresses, selectors []string, from, to uint64) ([]*model.Call, error) {
	res := make([]*model.Call, 0)
	for _, addr := range addresses {
		page := 1
		for {
			if page*DefaultPageSize > 10000 {
				return nil, utils.ErrTooManyRecords
			}
			o := Option{
				Page:       page,
				PageSize:   DefaultPageSize,
				StartBlock: int64(from),
				EndBlock:   int64(to),
				Asc:        true,
			}
			normalTxs, err := p.GetTransactions(addr, o)
			if err != nil {
				return nil, err
			}
			for _, t := range normalTxs {
				if t.To != addr || t.Error != "0" {
					continue
				}
				if len(selectors) != 0 && !utils.IsTargetCall(t.Input, selectors) {
					continue
				}
				bigVal, _ := new(big.Int).SetString(t.Value, 10)
				res = append(res, &model.Call{
					Number: utils.ParseStrToUint64(t.BlockNumber),
					Ts:     time.Unix(int64(utils.ParseStrToUint64(t.Timestamp)), 0),
					Index:  utils.ParseStrToUint64(t.Index),
					Hash:   t.Hash,
					From:   t.From,
					To:     t.To,
					Input:  t.Input,
					Value:  bigVal,
				})
			}
			if len(normalTxs) < DefaultPageSize {
				break
			}
			page++
		}

		hashToId := make(map[string]uint64)
		page = 1
		for {
			if page*DefaultPageSize > 10000 {
				return nil, utils.ErrTooManyRecords
			}
			o := Option{
				Page:       page,
				PageSize:   DefaultPageSize,
				StartBlock: int64(from),
				EndBlock:   int64(to),
				Asc:        true,
			}
			intTxs, err := p.GetInternalTransactions(addr, o)
			if err != nil {
				return nil, err
			}
			for _, t := range intTxs {
				if t.To != addr || t.Error != "0" || t.Type != "call" {
					continue
				}
				if len(selectors) != 0 && !utils.IsTargetCall(t.Input, selectors) {
					continue
				}
				if _, ok := hashToId[t.Hash]; !ok {
					hashToId[t.Hash] = 1
				} else {
					hashToId[t.Hash] += 1
				}
				bigVal, _ := new(big.Int).SetString(t.Value, 10)
				res = append(res, &model.Call{
					Number: utils.ParseStrToUint64(t.BlockNumber),
					Ts:     time.Unix(int64(utils.ParseStrToUint64(t.Timestamp)), 0),
					Hash:   t.Hash,
					Id:     hashToId[t.Hash],
					From:   t.From,
					To:     t.To,
					Input:  t.Input,
					Value:  bigVal,
				})
			}
			if len(intTxs) < DefaultPageSize {
				break
			}
			page++
		}

	}
	return res, nil
}

func (p *EtherscanProvider) GetContractInfo(address string) (*model.ContractInfo, error) {
	normal, err := p.GetFirstTransaction(address)
	var ret *model.ContractInfo
	if err != nil {
		return ret, err
	}
	var n uint64
	if normal != nil {
		n = utils.ParseStrToUint64(normal.BlockNumber)
		ret = &model.ContractInfo{
			Number:   n,
			Ts:       time.Unix(int64(utils.ParseStrToUint64(normal.Timestamp)), 0),
			Deployer: normal.From,
			Hash:     normal.Hash,
		}
	}
	internal, err := p.GetFirstInternalTransaction(address)
	if err != nil {
		return ret, err
	}
	if internal != nil {
		num := utils.ParseStrToUint64(internal.BlockNumber)
		addr := &internal.From
		if ret == nil || (addr != nil && num < n) {
			ret = &model.ContractInfo{
				Number:   num,
				Ts:       time.Unix(int64(utils.ParseStrToUint64(internal.Timestamp)), 0),
				Deployer: internal.From,
				Hash:     internal.Hash,
			}
		}
	}
	return ret, err
}

func (p *EtherscanProvider) GetContractFirstInvocation(address string) (ret uint64, err error) {
	normal, err := p.GetFirstTransaction(address)
	if normal != nil {
		ret = utils.ParseStrToUint64(normal.BlockNumber)
	}
	internal, err := p.GetFirstInternalTransaction(address)
	if err != nil {
		return
	}
	if internal != nil {
		num := utils.ParseStrToUint64(internal.BlockNumber)
		if ret == 0 {
			ret = num
		} else if num != 0 && num < ret {
			ret = num
		}
	}
	return
}

func (p *EtherscanProvider) GetFirstTransaction(address string) (*NormalTx, error) {
	txs, err := p.GetTransactions(address, Option{
		Page:       1,
		PageSize:   1,
		Asc:        true,
		StartBlock: 0,
		EndBlock:   MaxBlockNumber,
	})
	if err != nil {
		return nil, err
	}

	if len(txs) > 0 {
		return txs[0], nil
	}

	return nil, nil
}

func (p *EtherscanProvider) GetFirstInternalTransaction(address string) (*InternalTx, error) {
	txs, err := p.GetInternalTransactions(address, Option{
		Page:       1,
		PageSize:   1,
		Asc:        true,
		StartBlock: 0,
		EndBlock:   MaxBlockNumber,
	})
	if err != nil {
		return nil, err
	}
	if len(txs) > 0 {
		return txs[0], nil
	}

	return nil, nil
}

func (p *EtherscanProvider) GetTransactions(address string, o Option) ([]*NormalTx, error) {
	p.limiter.Wait(context.Background())
	k := p.nextKey()
	url := fmt.Sprintf(normalTxApi, p.baseUrl, strings.ToLower(address), o.StartBlock, o.EndBlock, o.Page, o.PageSize, toSortStr(o.Asc), k)
	ret, err := doFetchData[[]*NormalTx](url, p.proxy)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "invalid api key") {
		// p.deleteKey(k)
		return nil, fmt.Errorf("%v, %v", utils.ErrInvalidKey, k)
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "rate limit") {
		return nil, utils.ErrEtherscanRateLimit
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "smaller result dataset") {
		return nil, utils.ErrTooManyRecords
	}
	return ret, err
}

func (p *EtherscanProvider) GetInternalTransactions(address string, o Option) ([]*InternalTx, error) {
	p.limiter.Wait(context.Background())
	k := p.nextKey()
	url := fmt.Sprintf(internalTxApi, p.baseUrl, strings.ToLower(address), o.StartBlock, o.EndBlock, o.Page, o.PageSize, toSortStr(o.Asc), k)
	ret, err := doFetchData[[]*InternalTx](url, p.proxy)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "invalid api key") {
		// p.deleteKey(k)
		return nil, fmt.Errorf("%v, %v", utils.ErrInvalidKey, k)
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "rate limit") {
		return nil, utils.ErrEtherscanRateLimit
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "smaller result dataset") {
		return nil, utils.ErrTooManyRecords
	}
	return ret, err
}

func (p *EtherscanProvider) getLogs(topics0 string, o Option) ([]*EtherscanEvent, error) {
	p.limiter.Wait(context.Background())
	k := p.nextKey()
	url := fmt.Sprintf(logWithTopicsApi, p.baseUrl, o.StartBlock, o.EndBlock, topics0, o.Page, o.PageSize, k)
	if topics0 == "" {
		url = fmt.Sprintf(logApi, p.baseUrl, o.StartBlock, o.EndBlock, o.Page, o.PageSize, k)
	}
	ret, err := doFetchData[[]*EtherscanEvent](url, p.proxy)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "invalid api key") {
		// p.deleteKey(k)
		return nil, fmt.Errorf("%v, %v", utils.ErrInvalidKey, k)
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "rate limit") {
		return nil, utils.ErrEtherscanRateLimit
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "smaller result dataset") {
		return nil, utils.ErrTooManyRecords
	}
	return ret, err
}

func (p *EtherscanProvider) nextKey() string {
	p.l.Lock()
	defer p.l.Unlock()
	p.keyIter %= uint(len(p.apiKeys))
	key := p.apiKeys[p.keyIter]
	p.keyIter = (p.keyIter + 1) % uint(len(p.apiKeys))
	return key
}

func (p *EtherscanProvider) deleteKey(key string) {
	p.l.Lock()
	defer p.l.Unlock()
	p.apiKeys = utils.DeleteSliceElementByValue(p.apiKeys, key)
}

func doFetchData[T any](url, proxy string) (r T, err error) {
	log.Debug("invoke etherscan", "url", url)
	var resp EtherscanResponse[T]
	if err = utils.HttpGetObjectWithProxy(url, proxy, &resp); err != nil {
		err = fmt.Errorf("%v, url: %v", err, url)
		// log.Error("http get failed", "err", err, "url", url)
		return
	}
	if resp.Status != "1" && (resp.Message != noTransactionsFound && resp.Message != noRecordsFound) {
		err = fmt.Errorf("etherscan not ok: %s", resp.Message)
		// log.Error("etherscan get result falied", "err", err, "url", url)
		return
	}

	return resp.Result, nil
}

func toSortStr(asc bool) string {
	if asc {
		return "asc"
	} else {
		return "desc"
	}
}
