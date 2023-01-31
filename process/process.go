package processor

import (
	"app/cross_chain/anyswap"
	"app/dao"
	"app/model"
	"app/provider"
	"app/provider/chainbase"
	"app/svc"
	"app/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"sort"
)

type Processor struct {
	svc      *svc.ServiceContext
	chain    string
	provider *provider.Provider
}

func NewProcessor(svc *svc.ServiceContext, chain string) *Processor {
	p := svc.Providers.Get(chain)
	if p == nil {
		panic(fmt.Sprintf("%v: invalid provider", chain))
	}
	return &Processor{
		svc:      svc,
		chain:    chain,
		provider: p,
	}
}

const TIME_LAYOUT = "2006-01-02 15:04:05"

func (a *Processor) ProcessHopMultiMatched(d *dao.Dao, project_name string) error {
	var id_list []uint64

	//首先获取所有multi match的id
	stmt := "with t as (select * from " + d.Table() + " where direction = 'in' and project = $1 and match_id is not null)" +
		" select a.id from t a inner join t b on a.match_id = b.match_id and a.id != b.id group by a.id"
	err := d.DB().Select(&(id_list), stmt, project_name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(len(id_list))

	var visited = make(map[uint64]bool)
	for _, id := range id_list {
		if visited[id] { //如果已经处理过了，就直接处理下一批
			continue
		}

		da, _ := d.GetOne(id, "", -1)
		ele := da[0]
		in := d.GetDataWithTag("in", project_name, ele.MatchTag)
		out := d.GetDataWithTag("out", project_name, ele.MatchTag)

		for _, e := range in {
			visited[e.Id] = true
			chain_in, _ := e.ToChainId.Value()
			amount_in, _ := e.Amount.Value()

			for _, ee := range out {
				chain_out, _ := ee.ToChainId.Value()
				amount_out, _ := ee.Amount.Value()

				if !visited[ee.Id] && amount_in == amount_out && chain_in == chain_out {
					diff := e.Ts.Sub(ee.Ts)
					if diff.Hours() == 0 {
						d.UpdateMatchId(e.Id, ee.Id)
						d.UpdateMatchId(ee.Id, e.Id)
						visited[ee.Id] = true
						break
					}
				}
			}
		}
	}
	return nil
}

func (a *Processor) ProcessAcross(d *dao.Dao) error {
	var err error
	var pair []*uint64

	//已经test过，direction='out'的，没有重复match的in

	stmt := "select a.match_id from across a inner join across b on a.match_id = b.match_id and a.match_id is not null and a.id != b.id group by a.id"
	err = d.DB().Select(&pair, stmt)
	if err != nil {
		return err
	}

	var res []*string
	cnt := 0
	for _, p := range pair {
		//正常情况下只有一个fill_amount比较大
		stmt = fmt.Sprintf("select fill_amount from across where match_id = %d", *p)
		err = d.DB().Select(&res, stmt)
		if err != nil {
			return err
		}

		safe := 0
		for _, e := range res {
			if *e > "100" {
				safe++
			}
		}
		if safe > 1 {
			stmt = fmt.Sprintf("update across set safe='F' where match_id=%d", *p)
			_, err = d.DB().Exec(stmt)
			cnt++
		}
	}
	println(cnt)
	return err
}

//获取token信息

func (a *Processor) ProcessWithToken() error {
	for _, project := range a.svc.Config.Projects {
		if project == "Across" {
			continue
		}
		println(project)
		var infos model.ContractInfos
		token_map := make(map[string][]*string)

		data, err := a.svc.Dao.GetTokenChains(project)
		if len(data) == 0 || err != nil {
			return err
		}

		for _, e := range data {
			token_map[e.Chain] = append(token_map[e.Chain], &e.Token)
		}

		old_tokens, err := a.svc.Dao.GetContractInfos(project)
		if err != nil {
			return err
		}
		token_exist := make(map[string]bool)
		for _, e := range old_tokens {
			token_exist[e.Address] = true
		}

		for chain, tokens := range token_map {
			a.provider = a.svc.Providers.Get(chain)

			for _, token := range tokens {
				if _, ok := token_exist[*token]; ok {
					continue
				}

				info, err := a.provider.GetContractInfo(*token)
				if err != nil || info == nil {
					println(project, *token, chain)
					continue
				}

				info.Chain = chain
				info.Address = *token
				info.Project = project
				info.Type = "token"
				info.Safe = "T"
				infos = append(infos, info)

				if len(infos) == 20 {
					err = a.svc.Dao.SaveContractInfo(infos)
					if err != nil {
						return err
					}
					infos = *new(model.ContractInfos)
				}
			}
		}
		if err != nil {
			return err
		}

		if len(infos) != 0 {
			err = a.svc.Dao.SaveContractInfo(infos)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//通过判断deployer是否有label，进一步进行alert

func (a *Processor) AlertWithToken() error {
	/*err := a.ProcessWithToken()
	if err != nil {
		return err
	}*/

	for _, project := range a.svc.Config.Projects {
		if project == "Across" {
			continue
		}
		println(project)
		old_tokens, err := a.svc.Dao.GetContractInfos(project)
		if err != nil {
			return err
		}
		for _, token_info := range old_tokens {
			label, err := a.svc.Dao.QueryLabel(token_info.Address)
			if err != nil {
				fmt.Println(token_info.Address)
				return err
			}
			if label == false {
				label, err = a.svc.Dao.QueryLabel(token_info.Deployer)
				if err != nil {
					fmt.Println("deployer: ", token_info.Deployer)
					return err
				}
			}
			if label == false {
				err := a.svc.Dao.MarkUnsafeWithToken(token_info.Chain, token_info.Address)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *Processor) AddAnyswap() (int, error) {
	res, err := a.svc.ProjectsDao.GetUnmatchedAnyswap()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var results model.Results
	totalFetched := 0
	any := &anyswap.Anyswap{}
	topics0 := any.Topics0(a.chain)
	for _, r := range res {
		events, err := a.provider.GetLogsWithHash(r.Chain, r.Hash, topics0)
		if err != nil {
			return 0, err
		}
		totalFetched = len(events)
		if totalFetched == 0 {
			fmt.Println("unfetched hash: ", r.Chain, r.Hash)
			continue
		}
		sort.Sort(events)
		ret := any.Extract(r.Chain, events)
		if len(ret) == 0 {
			fmt.Println("CANNOT extract hash: ", r.Chain, r.Hash)
		}
		results = append(results, ret...)
	}
	err = a.svc.ProjectsDao.SaveAnyswap(results)
	totalFetched = len(results)
	if err != nil {
		return 0, fmt.Errorf("result save failed: %v", err)
	}
	return totalFetched, nil
}

func (a *Processor) DealMultiMatched() (int, error) {
	src_hashs, err := a.svc.ProjectsDao.GetMultiMatchedAnyswapSrcHash()
	if err != nil {
		return 0, err
	}

	cnt := 0
	for _, src_hash := range src_hashs {
		var in []*uint64
		var out []*uint64
		stmt := fmt.Sprintf("select id from aanyswap where tx_hash = '%s' order by block_number, log_index asc", src_hash.Hash)
		log.Debug(stmt)
		err := a.svc.ProjectsDao.DB().Select(&out, stmt)
		if err != nil {
			return 0, err
		}
		stmt = fmt.Sprintf("select id from aanyswap where src_tx_hash = '%s' order by block_number, log_index asc", src_hash.Hash)
		err = a.svc.ProjectsDao.DB().Select(&in, stmt)
		if err != nil {
			return 0, err
		}

		//如果in和out的数量不一致，则有问题
		if len(in) != len(out) {
			fmt.Println("alert: src_tx_hash = ", src_hash.Chain, src_hash.Hash)
			continue
		}
		for i := 0; i < len(in); i++ {
			err = a.svc.ProjectsDao.UpdateMatchId(*in[i], *out[i])
			if err != nil {
				return 0, err
			}
			err = a.svc.ProjectsDao.UpdateMatchId(*out[i], *in[i])
			if err != nil {
				return 0, err
			}
		}
		cnt++
	}
	return cnt, err
}

func (a *Processor) MarkWithToken() (int, error) {
	project := "Anyswap"
	anyUrl := "https://bridgeapi.multichain.org/v4/tokenlistv4/all"
	data, err := a.GetSupportedTokens(anyUrl)
	if err != nil {
		return 0, err
	}
	old_tokens, err := a.svc.Dao.GetContractInfos(project)
	if err != nil {
		return 0, err
	}
	for _, e := range old_tokens {
		chainId := utils.GetChainId(e.Chain)
		fmt.Println(chainId)
		token := fmt.Sprintf("evm" + e.Address)
		_, ok := data[token]
		if !ok { //如果在页面中查得到

		}
	}
	return 0, err
}

/*
@title: CheckAnyswapWithToken
@description: 通过contracts数据库里面标记了unsafe的token，检查有没有完成match的tx；或者标记了safe的token有没有unmatch的tx
*/

func (a *Processor) CheckAnyswapWithToken() (int, error) {
	project := "Anyswap"
	cnt := 0
	old_tokens, err := a.svc.Dao.GetContractInfos(project)
	if err != nil {
		return 0, err
	}
	for _, e := range old_tokens {
		if e.Deployer == "0xfa9da51631268a30ec3ddd1ccbf46c65fad99251" {
			continue
		}

		data, err := a.svc.ProjectsDao.GetAnyWithToken(e.Address)
		if err != nil {
			return cnt, err
		}
		matched := false
		for _, d := range data {
			if d.MatchId.Valid {
				matched = true
				break
			}
		}

		if matched == true && e.Safe == "F" {
			fmt.Println("unsafe token but matched! token address: ", e.Address)
			cnt++
		}
		if matched == false && e.Safe == "T" {
			fmt.Println("safe token but unmatched! token address: ", e.Address)
			cnt++
		}
	}
	return cnt, nil
}

//需要修改返回值类型

func (a *Processor) GetSupportedTokens(url string) (map[string]interface{}, error) {
	info, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(info))
	decoder.UseNumber()
	buffer := make(map[string]interface{})
	err = decoder.Decode(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer, err
}

const (
	interval  = 30 * 60
	batchSize = uint64(25000)
)

func (m *Processor) StartUpdateSy() {
	st := "select max(id) from synapse"
	var latest uint64
	err := m.svc.ProjectsDao.DB().Get(&latest, st)
	if err != nil {
		log.Error("get latest id failed", "err", err)
		return
	}
	var last uint64

	for last = 500001; last < latest; last = utils.Min(latest, last+batchSize) + 1 {
		go m.BeginUpdateSy(last, latest)
	}
}

func (m *Processor) BeginUpdateSy(last, latest uint64) {
	cnt := 0
	var err error
	var stmt string
	stmt = fmt.Sprintf("select %s from synapse where id>=$1 and id<=$2 and chain != 'cronos' and chain != 'boba' and from_address = ''", dao.Rows)
	var res = []*model.Data{}
	to := utils.Min(latest, last+batchSize)
	from := last
	err = m.svc.ProjectsDao.DB().Select(&res, stmt, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}
	println("from ", from, " to ", from+uint64(len(res)), " begins")

	var traces []*chainbase.SyChainbaseInfo
	for i, r := range res {
		chain := r.Chain
		hash := r.Hash
		trace, err := m.provider.GetTraces(chain, hash)
		if err != nil || len(trace) == 0 {
			cnt++
			fmt.Println(err)
			continue
		}
		trace[0].Hash = hash
		traces = append(traces, trace[0])

		if i%500 == 0 {
			err = m.svc.ProjectsDao.UpdateSy(traces)
			if err != nil {
				return
			}
			println("done ", to)
			traces = []*chainbase.SyChainbaseInfo{}
		}
	}
	println("from ", from, " to ", to, " unfetched: ", cnt)

	err = m.svc.ProjectsDao.UpdateSy(traces)
	if err != nil {
		return
	}
	return
}

/*
已地址为中心，查该地址in / out 数量不相等的情况
*/

func (a *Processor) GetUnmatchAddress(d *dao.Dao, table_name string) {
	stmt := fmt.Sprintf("with t as  (select count(direction), direction , to_address from %s group by direction, to_address order by to_address) "+
		"select a.to_address from t a inner join t b on a.to_address = b.to_address and a.direction='in' and b.direction='out' and a.count != b.count", table_name)
	log.Debug(stmt)
	var data []*string

	err := d.DB().Select(&data, stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (m *Processor) MarkTxWithFakeToken(d *dao.Dao, project string) {
	stmt := fmt.Sprintf("select address from contracts where project='%s' and safe = 'F'", project)
	token_chains := []*string{}
	err := d.DB().Select(&token_chains, stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, e := range token_chains {
		stmt = fmt.Sprintf("update %s set isFakeToken=1 where token='%s'", project, *e)
		_, err := d.DB().Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
	}
}
