package chainbase

import (
	"app/model"
	"app/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/time/rate"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	ChainbaseUrl = "https://api.chainbase.online/v1/dw/query"
)

var limiter *rate.Limiter

func SetupLimit(l int) {
	limiter = rate.NewLimiter(rate.Limit(l), 1)
}

type Provider struct {
	table           string
	apiKey          string
	enableTraceCall bool
	proxy           string
}

func NewProvider(table, apiKey string, enableTraceCall bool, proxy string) *Provider {
	return &Provider{
		table:           table,
		apiKey:          apiKey,
		enableTraceCall: enableTraceCall,
		proxy:           proxy,
	}
}

func (p *Provider) GetLatestNumber() (uint64, error) {
	limiter.Wait(context.Background())
	stmt := fmt.Sprintf("select max(number) as number from %v.blocks", p.table)
	log.Debug(stmt)
	ret, err := Exec[*Number](stmt, p.apiKey, p.proxy)
	if err != nil {
		return 0, err
	}
	if len(ret) == 0 {
		return 0, nil
	}
	return utils.ParseStrToUint64(ret[0].Number), nil
}

func (p *Provider) GetContractFirstCreatedNumber(address string) (uint64, error) {
	limiter.Wait(context.Background())
	if n, err := p.getContractCreatedNumber(address); err == nil && n > 0 {
		return n, nil
	}
	if p.enableTraceCall {
		return p.getAddressFirstCallByTraceCall(address)
	}
	return p.getAddressFirstCallByTransaction(address)
}

func (p *Provider) getContractCreatedNumber(address string) (uint64, error) {
	stmt := fmt.Sprintf("select block_number as number from %v.contracts where contract_address = '%v'", p.table, address)
	log.Debug(stmt)
	ret, err := Exec[*Number](stmt, p.apiKey, p.proxy)
	if err != nil {
		return 0, err
	}
	if len(ret) == 0 {
		return 0, nil
	}
	return utils.ParseStrToUint64(ret[0].Number), nil
}

func (p *Provider) getAddressFirstCallByTransaction(address string) (uint64, error) {
	var n1, n2 uint64
	stmt1 := fmt.Sprintf("select min(block_number) as number from %v.transactions where contract_address = '%v' and status = 1", p.table, address)
	log.Debug(stmt1)
	ret1, err := Exec[*Number](stmt1, p.apiKey, p.proxy)
	if err != nil {
		return 0, err
	}
	if len(ret1) > 0 {
		n1 = utils.ParseStrToUint64(ret1[0].Number)
	}

	stmt2 := fmt.Sprintf("select min(block_number) as number from %v.transactions where to_address = '%v' and status = 1", p.table, address)
	log.Debug(stmt2)
	ret2, err := Exec[*Number](stmt2, p.apiKey, p.proxy)
	if err != nil {
		return 0, err
	}
	if len(ret2) > 0 {
		n2 = utils.ParseStrToUint64(ret2[0].Number)
	}
	if n1 == 0 && n2 == 0 {
		return 0, nil
	} else if n1 == 0 {
		return n2, nil
	} else if n2 == 0 {
		return n1, nil
	}
	return utils.Min(n1, n2), nil
}

func (p *Provider) getAddressFirstCallByTraceCall(address string) (uint64, error) {
	stmt := fmt.Sprintf("select min(block_number) as number from %v.trace_calls where to_address = '%v' and error = ''", p.table, address)
	log.Debug(stmt)
	ret, err := Exec[*Number](stmt, p.apiKey, p.proxy)
	if err != nil {
		return 0, err
	}
	if len(ret) > 0 {
		return utils.ParseStrToUint64(ret[0].Number), nil
	}
	return 0, nil
}

func (p *Provider) GetLogs(addresses []string, topics0 []string, from, to uint64) (model.Events, error) {
	limiter.Wait(context.Background())
	if len(addresses) == 0 {
		return nil, nil
	}
	res := make(model.Events, 0)

	stmt := fmt.Sprintf("select * from %v.transaction_logs where block_number >= %v and block_number <= %v and %v", p.table, from, to, formatOrCondition("address", addresses))
	if len(topics0) > 0 {
		stmt += " and " + formatOrCondition("topic0", topics0)
	}
	log.Debug(stmt)
	ret, err := Exec[*Log](stmt, p.apiKey, p.proxy)
	if err != nil {
		return nil, err
	}
	ret = filterDuplicateLogs(ret)
	for _, l := range ret {
		ts, err := utils.ParseDateTime(l.Ts)
		if err != nil {
			log.Error("invalid ts from chainbase", "chain", p.table, "block", l.Number, "err", err)
		}
		topics := make([]string, 0)
		for i := 0; i < int(l.TopicsCnt); i++ {
			switch i {
			case 0:
				topics = append(topics, l.T0)
			case 1:
				topics = append(topics, l.T1)
			case 2:
				topics = append(topics, l.T2)
			case 3:
				topics = append(topics, l.T3)
			}
		}
		res = append(res, &model.Event{
			Number:  utils.ParseStrToUint64(l.Number),
			Ts:      ts,
			Index:   l.Index,
			Hash:    l.Hash,
			Id:      l.LogIndex,
			Address: l.Address,
			Topics:  topics,
			Data:    l.Data,
		})
	}
	return res, nil
}

func (p *Provider) GetCalls(addresses, selectors []string, from, to uint64) ([]*model.Call, error) {
	limiter.Wait(context.Background())
	if len(addresses) == 0 {
		return nil, nil
	}
	res := make([]*model.Call, 0)
	stmt := fmt.Sprintf("select * from %v.transactions where block_number >= %v and block_number <= %v and status = 1 and %v", p.table, from, to, formatOrCondition("to_address", addresses))
	if p.enableTraceCall {
		stmt = fmt.Sprintf("select * from %v.trace_calls where block_number >= %v and block_number <= %v and error = '' and %v and call_type = 'call'", p.table, from, to, formatOrCondition("to_address", addresses))
		if len(selectors) != 0 {
			stmt += fmt.Sprintf(" and %v", formatOrCondition("method_id", trimPrefix(selectors, "0x")))
		}
	}
	log.Debug(stmt)
	ret, err := Exec[*Trace](stmt, p.apiKey, p.proxy)
	ret = filterDuplicateCalls(ret)
	if err != nil {
		return nil, err
	}
	sort.Stable(Traces(ret))
	id := uint64(0)
	prevHash := ""
	for _, t := range ret {
		if prevHash == "" || prevHash != t.Hash {
			// next is internal tx
			if len(t.TraceAddress) != 0 {
				id = 1
			} else {
				//next is external tx
				id = 0
			}
		} else {
			id += 1
		}
		prevHash = t.Hash
		if len(selectors) != 0 && !utils.IsTargetCall(t.Input, selectors) {
			continue
		}
		bigVal, _ := new(big.Int).SetString(t.Value, 10)
		ts, err := utils.ParseDateTime(t.Ts)
		if err != nil {
			log.Error("invalid ts from chainbase", "chain", p.table, "block", t.Number, "err", err)
		}
		res = append(res, &model.Call{
			Number: utils.ParseStrToUint64(t.Number),
			Ts:     ts,
			Index:  t.Index,
			Hash:   t.Hash,
			Id:     id,
			From:   t.From,
			To:     t.To,
			Input:  t.Input,
			Value:  bigVal,
		})
	}
	return res, nil
}

func Exec[T any](stmt, apiKey, proxy string) ([]T, error) {
	res := make([]T, 0)
	taskId := ""
	page := uint(0)
	for {
		ret, err := exec[T](stmt, taskId, apiKey, proxy, page)
		if err != nil {
			return nil, err
		}
		if ret.Message != "ok" {
			return nil, fmt.Errorf("chainbase error: %v", ret.Message)
		}
		if ret.Data.ErrMsg != "" {
			return nil, fmt.Errorf("chainbase error: %v", ret.Data.ErrMsg)
		}
		res = append(res, ret.Data.Result...)
		if ret.Data.NextPage != 0 {
			taskId = ret.Data.TaskId
			page = ret.Data.NextPage
		} else {
			break
		}
	}
	return res, nil
}

func exec[T any](stmt, taskId, apiKey, proxy string, page uint) (ret *Result[T], err error) {
	ret = &Result[T]{}
	u, err := url.Parse(ChainbaseUrl)
	if err != nil {
		return nil, err
	}
	reqBody := map[string]any{
		"query": stmt,
	}
	if taskId != "" && page != 0 {
		reqBody["task_id"] = taskId
		reqBody["page"] = page
	}
	opt := utils.HttpOption{
		Method: http.MethodPost,
		Url:    u,
		Header: map[string]string{
			"x-api-key": apiKey,
		},
		RequestBody: reqBody,
		Response:    ret,
		Proxy:       proxy,
	}
	err = opt.Send(context.Background())
	return
}

func trimPrefix(ss []string, prefix string) []string {
	ret := make([]string, 0, len(ss))
	for _, s := range ss {
		ret = append(ret, strings.TrimPrefix(s, prefix))
	}
	return ret
}

func formatOrCondition(field string, args []string) string {
	if len(args) == 0 {
		return ""
	}
	cond := "("
	for idx, arg := range args {
		cond += fmt.Sprintf("%v = '%v'", field, arg)
		if idx < len(args)-1 {
			cond += " or "
		}
	}
	cond += ")"
	return cond
}

func filterDuplicateLogs(logs []*Log) []*Log {
	res := make([]*Log, 0, len(logs))
	set := make(map[string]struct{})
	for _, l := range logs {
		key := l.Number + "-" + strconv.FormatUint(l.Index, 10) + "-" + strconv.FormatUint(l.LogIndex, 10)
		if _, ok := set[key]; ok {
			continue
		}
		res = append(res, l)
		set[key] = struct{}{}
	}
	return res
}

func filterDuplicateCalls(calls []*Trace) []*Trace {
	res := make([]*Trace, 0, len(calls))
	set := make(map[string]struct{})
	for _, c := range calls {
		tmp, _ := json.Marshal(c.TraceAddress)
		key := c.Number + "-" + strconv.FormatUint(c.Index, 10) + "-" + string(tmp)
		if _, ok := set[key]; ok {
			continue
		}
		res = append(res, c)
		set[key] = struct{}{}
	}
	return res
}

func (p *Provider) GetContractInfo(token string) (*model.ContractInfo, error) {
	limiter.Wait(context.Background())

	stmt := fmt.Sprintf("select * from %v.contracts where contract_address = '%v'", p.table, token)
	log.Debug(stmt)
	ret, err := Exec[*Trace](stmt, p.apiKey, p.proxy)
	if err != nil || len(ret) == 0 {
		return nil, err
	}

	ts, err := utils.ParseDateTime(ret[0].Ts)
	res := &model.ContractInfo{
		Hash:     ret[0].Hash,
		Number:   utils.ParseStrToUint64(ret[0].Number),
		Ts:       ts,
		Deployer: ret[0].From,
	}

	return res, nil
}

func (p *Provider) GetLogsWithHash(chain, hash string, topics0 []string) (model.Events, error) {
	limiter.Wait(context.Background())
	res := make(model.Events, 0)

	stmt := fmt.Sprintf("select * from %v.transaction_logs where transaction_hash  = '%s'", chain, hash)
	if len(topics0) > 0 {
		stmt += " and " + formatOrCondition("topic0", topics0)
	}
	log.Debug(stmt)
	ret, err := Exec[*Log](stmt, p.apiKey, p.proxy)
	if err != nil {
		return nil, err
	}
	ret = filterDuplicateLogs(ret)
	for _, l := range ret {
		ts, err := utils.ParseDateTime(l.Ts)
		if err != nil {
			log.Error("invalid ts from chainbase", "chain", p.table, "block", l.Number, "err", err)
		}
		topics := make([]string, 0)
		for i := 0; i < int(l.TopicsCnt); i++ {
			switch i {
			case 0:
				topics = append(topics, l.T0)
			case 1:
				topics = append(topics, l.T1)
			case 2:
				topics = append(topics, l.T2)
			case 3:
				topics = append(topics, l.T3)
			}
		}
		res = append(res, &model.Event{
			Number:  utils.ParseStrToUint64(l.Number),
			Ts:      ts,
			Index:   l.Index,
			Hash:    l.Hash,
			Id:      l.LogIndex,
			Address: l.Address,
			Topics:  topics,
			Data:    l.Data,
		})
	}
	return res, nil
}

type SyChainbaseInfo struct {
	Hash string
	Ts   string `json:"block_timestamp"`
	From string `json:"from_address"`
}

func (p *Provider) GetTraces(chain, hash string) ([]*SyChainbaseInfo, error) {
	//limiter.Wait(context.Background())
	stmt := fmt.Sprintf("select block_timestamp, from_address from %s.transactions where transaction_hash = '%s'", chain, hash)
	log.Debug(stmt)
	ret, err := Exec[*SyChainbaseInfo](stmt, p.apiKey, p.proxy)
	return ret, err
}

func (p *Provider) GetSender(chain, hash string) (string, error) {
	type S struct {
		Sender string `json:"from_address"`
	}
	stmt := fmt.Sprintf("select from_address from %s.transactions where transaction_hash='%s'", chain, hash)
	ret, err := Exec[*S](stmt, p.apiKey, p.proxy)
	if len(ret) == 0 || err != nil {
		log.Warn("GetSender(): failed to get sender", "Chain ", chain, "Hash ", hash, "Error ", err.Error())
		return "", err
	}
	return ret[0].Sender, err
}
