package dao

import (
	"app/model"
	"app/provider/chainbase"
	"app/utils"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stringx"
)

const (
	SQL_MAX_PLACEHOLDERS = 65535
)

var (
	resultInsertFieldNames = builder.RawFieldNames(&model.Data{}, true)
	resultInsertRows       = strings.Join(stringx.Remove(resultInsertFieldNames, "id", "match_id", "profit", "from_address_error", "to_address_profit", "token_profit_error", "isfaketoken"), ",")
	resultInsertTags       = strings.Join(slicesWithPrefix(stringx.Remove(resultInsertFieldNames, "id", "match_id", "profit", "from_address_error", "to_address_profit", "token_profit_error", "isfaketoken"), ":"), ",")

	resultUpdateFieldNames = []string{"match_id", "from_chain", "from_address", "to_chain", "to_address"}
	resultUpdateRows       = strings.Join(resultUpdateFieldNames, ",")
	resultUpdateTags       = builder.PostgreSqlJoin(resultUpdateFieldNames)
	// resultUpdateWithPlaceholders = builder.PostgreSqlJoin(resultUpdateFieldNames)
	// resultWithPlaceholder = builder.PostgreSqlJoin(stringx.Remove(resultFieldNames, "id"))

	contractInsertFieldNames = builder.RawFieldNames(&model.ContractInfo{}, true)
	contractInsertRows       = strings.Join(stringx.Remove(contractInsertFieldNames, "id"), ",")
	contractInsertTags       = strings.Join(slicesWithPrefix(stringx.Remove(contractInsertFieldNames, "id"), ":"), ",")
)

type Dao struct {
	db         *sqlx.DB
	table      string
	otherTable map[string]string
}

func NewDao(host string) *Dao {
	db, err := sqlx.Connect("postgres", host)
	if err != nil {
		panic(err)
	}
	return &Dao{
		db:    db,
		table: "common_cross_chain",
		otherTable: map[string]string{
			"contracts": "contracts",
			"labels":    "multi_chain_addresses",
		},
	}
}

func NewAnyDao(host string) *Dao {
	db, err := sqlx.Connect("postgres", host)
	if err != nil {
		panic(err)
	}
	return &Dao{
		db:    db,
		table: "anyswap",
		otherTable: map[string]string{
			"synapse": "synapse",
		},
	}
}

func (d *Dao) LatestId(project string) (latest uint64, err error) {
	stmt := fmt.Sprintf("select max(id) from %s", project)
	err = d.db.Get(&latest, stmt)
	return
}

func (d *Dao) LastMatchedId(project string, Id uint64) (last uint64, err error) {
	stmt := fmt.Sprintf("select max(id) from %s where id >= %d and direction = 'in' and match_id is null", project, Id)
	err = d.db.Get(&last, stmt)

	return
}

func (d *Dao) SaveContractInfo(info model.ContractInfos) error {
	if len(info) == 0 {
		return nil
	}
	stmt := fmt.Sprintf("insert into %s (%s) values (%s)", d.otherTable["contracts"], contractInsertRows, contractInsertTags)
	_, err := d.db.NamedExec(stmt, info)
	return err
}

func (d *Dao) GetContractInfos(project string) (model.ContractInfos, error) {
	var res model.ContractInfos
	stmt := fmt.Sprintf("select * from %s where project = '%s'", d.otherTable["contracts"], project)
	err := d.db.Select(&res, stmt)
	if err != nil {
		return nil, err
	}
	fmt.Println("existed tokens: ", len(res))
	return res, nil
}

func (d *Dao) Save(results model.Datas, table ...string) (err error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error("save results rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	maxInsert := SQL_MAX_PLACEHOLDERS / 17
	for i := 0; i < len(results); i += maxInsert {
		batch := results[i:utils.Min(i+maxInsert, len(results))]
		err = d.save(tx, batch, table...)
		if err != nil {
			return
		}
	}
	return nil
}

func (d *Dao) save(tx *sqlx.Tx, results model.Datas, table ...string) error {
	if len(table) > 0 {
		d.table = table[0]
	}
	if len(results) == 0 {
		return nil
	}
	stmt := fmt.Sprintf("insert into %s (%s) values (%s)", d.table, resultInsertRows, resultInsertTags)
	fmt.Println(stmt)
	_, err := tx.NamedExec(stmt, results)
	return err
}

func (d *Dao) Update(project string, results model.Datas) (err error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error("save trades rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	stmt := fmt.Sprintf("update %s set %s where id = $1", project, resultUpdateTags)
	for _, r := range results {
		_, err = d.db.Exec(stmt, r.Id, r.MatchId, r.FromChainId, r.FromAddress, r.ToChainId, r.ToAddress)
		if err != nil {
			log.Error("update failed", "err", err)
		}
	}
	return
}

func (d *Dao) Delete(project string, id uint64) error {
	stmt := fmt.Sprintf("delete from %s where id = %d", project, id)
	_, err := d.db.Exec(stmt)
	return err
}

func (d *Dao) UpdateAnyswapMatchTag(project string, results model.Datas) int {
	tx, err := d.db.Beginx()
	cnt := 0
	if err != nil {
		return cnt
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error("save trades rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	filedName := []string{"match_tag"}
	FiledName := builder.PostgreSqlJoin(filedName)
	stmt := fmt.Sprintf("update %s set %s where id = $1", project, FiledName)
	for _, r := range results {
		_, err = d.db.Exec(stmt, r.Id, r.MatchTag)
		if err != nil {
			log.Error("update failed", "err", err)
		} else {
			cnt++
		}
	}
	return cnt
}

func (d *Dao) UpdateMatchId(id, match_id uint64) error {
	stmt := "UPDATE " + d.table + " SET match_id = $2 WHERE id = $1"
	_, err := d.db.Exec(stmt, id, match_id)
	return err
}

func (d *Dao) LastUpdate(chain, project string) (uint64, error) {
	var last uint64
	stmt := "select number from common_cross_chain where chain = $1 and project = $2 order by number desc limit 1"
	err := d.db.Get(&last, stmt, chain, project)
	if err == sql.ErrNoRows {
		err = nil
	}
	return last, err
}

// @title	SelectProject
// @description	Get data of one project with the limitation of direction and whether the match_id is empty
// @auth	Hu xiaohui
// @param	projectName	string	the project name
// @param
// @return	complete infomation of the data
func (d *Dao) SelectProject(projectName, direct, empty string) []model.Data {
	if len(projectName) == 0 || len(direct) == 0 {
		return nil
	}

	res := []model.Data{}

	var err error
	if empty == "false" {
		stmt := "SELECT * FROM common_cross_chain WHERE project = $1 AND direction = $2 AND match_id IS NOT NULL order by number desc"
		err = d.db.Select(&res, stmt, projectName, direct)
		//fmt.Println(*res[0])
	} else if empty == "true" {
		stmt := "SELECT * FROM common_cross_chain WHERE project = $1 AND direction = $2 AND match_id IS NULL order by number desc"
		err = d.db.Select(&res, stmt, projectName, direct)
	} else if empty == "all" {
		stmt := "SELECT * FROM common_cross_chain WHERE project = $1 AND direction = $2 "
		err = d.db.Select(&res, stmt, projectName, direct)
	} else {
		log.Warn("unavilable 'empty'")
		return nil
	}

	//res = append(res, r)

	if err != nil {
		fmt.Println("err:", err)
		return nil
	}
	return res
}

// @title	GetOne
// @description	Allows to get one data from pg with its (id, hash, matchId) OR just one info
// @auth	Hu xiaohui
// @param	id, hash, match_id
// @return	complete infomation of the data
func (d *Dao) GetOne(id uint64, hash string, matchId int64) (model.Datas, error) {
	var res model.Datas
	var res_ model.Datas
	var _res model.Datas

	var err error
	if id >= 0 {
		err = d.db.Select(&res, "SELECT * from common_cross_chain WHERE id = $1", id)
	}

	if len(hash) == 66 {
		err = d.db.Select(&res_, "SELECT * FROM common_cross_chain WHERE hash = $1", hash)
		if (len(res[0].Hash) != 0) && (res[0].Hash != res_[0].Hash) {
			log.Warn("paramters are not matched with each other")
			return nil, err
		}
		res = res_
	}

	if matchId >= 0 {
		err = d.db.Get(&_res, "SELECT * from common_cross_chain WHERE match_id = $1", matchId)
		if (len(res[0].Hash) != 0) && (res[0].Hash != _res[0].Hash) {
			log.Warn("paramters are not matched with each other")
			return nil, err
		}
		res = _res
	}

	return res, err
}

func (d *Dao) GetTokenChains(project_name string) (model.TokenChains, error) {
	var tokens model.TokenChains
	stmt := "select token, chain from common_cross_chain where project = $1 group by token, chain order by count(token) asc"
	err := d.db.Select(&tokens, stmt, project_name)
	if err != nil {
		return nil, err
	}
	println("all tokens: ", len(tokens))
	return tokens, err
}

func (d *Dao) QueryLabel(address string) (bool, error) {
	stmt := fmt.Sprintf("select name_tag, labels, contract_name, token_name from %s where address = '%s'", d.otherTable["contracts"], address)
	log.Debug(stmt)
	var res []*model.LabelInfo
	err := d.db.Select(&res, stmt)
	if err != nil || len(res) == 0 ||
		strings.Contains(res[0].Labels.String, "EXPLOIT") {
		return false, err
	}
	return true, err
}

func (d *Dao) MarkUnsafeWithToken(chain, token string) error {
	stmt := fmt.Sprintf("update %s set safe='F' where chain='%s' and address='%s'", d.otherTable["contracts"], chain, token)
	log.Debug(stmt)
	_, err := d.db.Exec(stmt)
	return err
}

func (d *Dao) GetDataWithTag(direction, project_name, match_tag string) []*model.Data {
	var res model.Datas
	stmt := "select * from " + d.table + " where direction=$1 and project=$2 and match_tag=$3 " +
		"and (to_chain_id=1 or to_chain_id=10 or to_chain_id=56 or to_chain_id=137 or to_chain_id=250 or to_chain_id=42161 or to_chain_id=43114) order by ts asc"
	err := d.db.Select(&res, stmt, direction, project_name, match_tag)
	if err != nil {
		log.Warn("failed to get data from pg: ", err)
	}
	return res
}

func (d *Dao) GetDataWithToken(project, token string) model.Datas {
	var res model.Datas
	stmt := fmt.Sprintf("select * from %s where project='%s' and token='%s'", d.table, project, token)
	err := d.db.Select(&res, stmt)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func (d *Dao) GetOneMatched(data model.Data) ([]model.Data, error) {
	var res []model.Data
	stmt := "select * from " + d.Table() + " where project = $1 and to_chain_id = $1 and match_tag = $3 and id != $4"
	err := d.db.Select(&res, stmt, data.ToChainId, data.MatchTag, data.Id)
	return res, err
}

func (d *Dao) GetBatchMatched(project_name string) ([]*model.MatchedId, error) {
	res := []*model.MatchedId{}

	stmt := "with t as (select * from " + d.table + " where project = $1 and match_id is null)" +
		" select t1.id as dest_id, t2.id as src_id from t t1 inner join t t2 " +
		"on t1.match_tag = t2.match_tag and t1.to_address = t2.to_address and t1.direction='in' and t2.direction='out' and t1.to_chain_id = t2.to_chain_id"
	err := d.db.Select(res, stmt, project_name)

	fmt.Println(len(res))
	return res, err
}

func (d *Dao) GetBatchMatched_Hop_Stargate(project_name string) ([]*model.MatchedId, error) {
	res := []*model.MatchedId{}

	stmt := "with t as (select * from " + d.table + " where project = $1 and match_id is null)" +
		" select t1.id as dest_id, t2.id as src_id from t t1 inner join t t2 " +
		"on t1.match_tag = t2.match_tag and t1.direction='in' and t2.direction='out' and t1.to_chain_id = t2.to_chain_id"
	err := d.db.Select(res, stmt, project_name)

	fmt.Println(len(res))
	return res, err
}

func (d *Dao) DB() *sqlx.DB { return d.db }

func (d *Dao) Table() string { return d.table }

func slicesWithPrefix(s []string, prefix string) []string {
	ret := make([]string, len(s))
	for i := range s {
		ret[i] = prefix + s[i]
	}
	return ret
}

func getSqlxNamedTagsForUpdate(s []string) []string {
	ret := make([]string, len(s))
	for i := range s {
		ret[i] = fmt.Sprintf(`%s=:%s`, s[i], s[i])
	}
	return ret
}

// ========================================================================================================================
// ========================================================================================================================
// ========================================================================================================================

var ChainName = map[string]string{
	"1":     "ethereum",
	"10":    "optimism",
	"56":    "bsc",
	"137":   "polygon",
	"250":   "fantom",
	"42161": "arbitrum",
	"43114": "avalanche",
}

var (
	FieldNames = builder.RawFieldNames(&model.Data{}, true)
	Rows       = strings.Join(stringx.Remove(FieldNames, "id"), ",")
	Tags       = strings.Join(slicesWithPrefix(stringx.Remove(FieldNames, "id"), ":"), ",")
)

type Res struct {
	Hash  string `db:"match_tag"`
	Chain string `db:"from_chain"`
}

func (d *Dao) GetUnmatchedAnyswap() ([]*Res, error) {
	stmt := "select match_tag, from_chain from aanyswap where direction = 'in' and (from_chain = '1' or from_chain = '10' or from_chain = '56' or from_chain ='137' or from_chain ='250' or from_chain ='42161' or from_chain = '43114') and (to_chain = '1' or to_chain = '10' or to_chain='56' or to_chain='137' or to_chain='250' or to_chain='42161' or to_chain='43114') and match_id is null"
	var res []*Res
	err := d.db.Select(&res, stmt)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		e.Chain = ChainName[e.Chain]
	}
	return res, err
}

func (d *Dao) GetMultiMatchedAnyswapSrcHash() ([]*Res, error) {
	stmt := "select a.match_tag, a.from_chain from aanyswap a inner join aanyswap a2 on a.match_id = a2.match_id and a.id != a2.id and a.direction = 'in' group by a.match_tag, a.from_chain"
	var res []*Res
	err := d.db.Select(&res, stmt)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		e.Chain = ChainName[e.Chain]
	}
	return res, err
}

func (d *Dao) GetAnyWithToken(token string) ([]*model.Data, error) {
	var res []*model.Data
	stmt := fmt.Sprintf("select * from %s where token = '%s'", d.table, token)
	log.Debug(stmt)
	err := d.db.Select(&res, stmt)
	return res, err
}

func (d *Dao) SaveAnyswap(results model.Datas) (err error) {
	var datas []*model.Data
	for _, r := range results {
		data := &model.Data{
			Id:          r.Id,
			Chain:       r.Chain,
			Number:      r.Number,
			TxIndex:     r.TxIndex,
			Hash:        r.Hash,
			LogIndex:    r.LogIndex,
			Contract:    r.Contract,
			Direction:   r.Direction,
			FromAddress: r.FromAddress,
			FromChainId: r.FromChainId,
			ToAddress:   r.ToAddress,
			ToChainId:   r.ToChainId,
			Token:       r.Token,
			Amount:      r.Amount,
			MatchTag:    r.MatchTag,
		}
		datas = append(datas, data)
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error("save results rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	maxInsert := SQL_MAX_PLACEHOLDERS / 17
	for i := 0; i < len(datas); i += maxInsert {
		batch := datas[i:utils.Min(i+maxInsert, len(results))]
		err = d.saveAny(tx, batch)
		if err != nil {
			return
		}
	}
	return nil
	return nil
}

func (d *Dao) saveAny(tx *sqlx.Tx, results []*model.Data) error {
	if len(results) == 0 {
		return nil
	}
	stmt := fmt.Sprintf("insert into aanyswap (%s) values (%s)", Rows, Tags)
	_, err := tx.NamedExec(stmt, results)
	return err
}

// ========================================================================================================================

func (d *Dao) GetSyData(stmt string) ([]*model.Data, error) {
	var res []*model.Data
	var err error
	if err = d.db.Select(&res, stmt); err == nil {
		return res, nil
	}
	return res, err
}

func (d *Dao) UpdateSy(results []*chainbase.SyChainbaseInfo) error {
	if len(results) == 0 {
		return nil
	}
	SyUpdateFieldNames := []string{"ts", "from_address"}
	SyUpdateTags := builder.PostgreSqlJoin(SyUpdateFieldNames)

	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error("save trades rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	stmt := fmt.Sprintf("update synapse set %s where hash = $1", SyUpdateTags)
	// fmt.Println(stmt)
	for _, r := range results {
		_, err = d.db.Exec(stmt, r.Hash, r.Ts, r.From)
		if err != nil {
			log.Error("update failed", "err", err)
		}
	}
	return err
}
