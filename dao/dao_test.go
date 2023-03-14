package dao

import (
	"app/model"
	"fmt"
	"reflect"
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
	//	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	d := NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")

	stmt := fmt.Sprintf("select id from %s where hash = $1 and log_index = $2", "Anyswap")
	var id uint64
	var hash = "0x01e04e7936aa24195a0beec29d2fbd6be518ae103284169e22766f75bdcf4084"
	log_index := 61
	err := d.db.Get(&id, stmt, hash, log_index)
	if err != nil || id == 0 {
		fmt.Println(err)
	}
	fmt.Println(id)
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

func TestNewAnyDao(t *testing.T) {
	stmt := "select * from anyswap where id = 77000000"
	d := NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	var da model.Data
	err := d.DB().Get(&(da), stmt)
	if err != nil {
		fmt.Println(err)
	}
	if reflect.DeepEqual(da, model.Data{}) {
		println("nil")
	}
}
