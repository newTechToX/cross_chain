package dao

import (
	"app/model"
	"app/utils"
	"fmt"
	"testing"
)

func TestDao(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	var f *uint64

	err := d.db.Get(&f, "select match_id from common_cross_chain where hash = '0x400b912ee6f55c80facf3e0f14347a1ad994fc241cd888dd00e31ec8db327915'")
	fmt.Println(err)
	fmt.Println(*f)
	fmt.Println(resultInsertRows)
	fmt.Println(resultInsertTags)
	fmt.Println(resultUpdateRows)
	fmt.Println(resultUpdateTags)
	// fmt.Println(fmt.Sprintf("insert into %s (%s) values (%s)", d.table, resultInsertRows, resultInsertTags))
	// fmt.Printf("update %s set %s where id = :id", d.table, resultUpdateTags)
}

func TestGet(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")

	stmt := "SELECT * FROM common_cross_chain WHERE id = $1"

	res := model.Data{}
	// _ = d.db.Get(&res, stmt, 1156329)
	// fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)

	// _ = d.db.Get(&res, stmt, 1131020)
	// fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)
	// _ = d.db.Get(&res, stmt, 1131019)
	// fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)
	err := d.db.Get(&res, stmt, 15159058)
	// fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)
	utils.PrintPretty(res)
	fmt.Println(err)
}

func TestUpdate(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "UPDATE common_cross_chain SET match_id = $2 WHERE id = $1"
	_, _ = d.db.Exec(stmt, 85608, 1001689)
	_, _ = d.db.Exec(stmt, 85623, 1001701)
	_, _ = d.db.Exec(stmt, 1001701, 85623)

}

func TestDelCol(t *testing.T) {
	d := NewAnyDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	println(d.table)
}

func TestDao_GetSynapseData(t *testing.T) {
	d := NewAnyDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := fmt.Sprintf("select %s from synapse where id = 1", Rows)
	res, err := d.GetSyData(stmt)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res[0])
}

func TestDa(t *testing.T) {
	d := NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "update synapse set ()"
	var r []*model.Data
	err := d.db.Select(&r, stmt)
	if err != nil {
		println(fmt.Println(err))
	}
}

func TestBigFloat_Cmp(t *testing.T) {
	stmt := fmt.Sprintf("select %s from across limit 10", model.ResultRows)
	d := NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	var data model.Datas
	err := d.DB().Select(&data, stmt)
	if err != nil {
		fmt.Println(err)
	} else {
		println(len(data))
	}
}
