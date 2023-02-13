package processor

import (
	"app/dao"
	"app/model"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// var d = dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
var d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")

func TestProcess_ProcessMultiMatched(t *testing.T) {
	a := Processor{}
	_ = a.ProcessHopMultiMatched(d, "Hop")
}

func Test_1(t *testing.T) {
	TIME_LAYOUT := "2006-01-02 15:04:05"
	ss := "2021-06-20 13:49:46.000"
	tt := "2021-06-17 13:21:36.000"
	sss, _ := time.Parse(TIME_LAYOUT, ss)
	ttt, _ := time.Parse(TIME_LAYOUT, tt)
	r := sss.Sub(ttt)

	if r.Hours() >= 1 {
		fmt.Println(sss.Sub(ttt).Hours())
	}
}

func TestProcess_Process_Across(t *testing.T) {
	a := Processor{}
	err := a.ProcessAcross(d)
	if err != nil {
		fmt.Println(err)
	}
}

func TestProcessor_AlertWithToken(t *testing.T) {
	a := Processor{}
	for _, p := range a.svc.Config.Projects {
		fmt.Println(p)
	}
}

func TestTokenApi(t *testing.T) {
	jsonData := []byte(`{"Name":"Eve","Age":6,"Parents":["Alice","Bob"]}`)

	var v interface{}
	json.Unmarshal(jsonData, &v)
	data := v.(map[string]interface{})

	for k, v := range data {
		switch v := v.(type) {
		case string:
			fmt.Println(k, v, "(string)")
		case float64:
			fmt.Println(k, v, "(float64)")
		case []interface{}:
			fmt.Println(k, "(array):")
			for i, u := range v {
				fmt.Println("    ", i, u)
			}
		default:
			fmt.Println(k, v, "(unknown)")
		}
	}
}

func TestProcessor_BeginUpdateSy(t *testing.T) {
	m := Processor{}
	m.BeginUpdateSy(1, 10)
}

func TestNewProcessor(t *testing.T) {
	d := dao.NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	//a := Processor{}
	stmt := "select match_id from across where safe is not null group by match_id"
	var res []*uint64
	err := d.DB().Select(&res, stmt)
	if err != nil {
		fmt.Println(err)
	}
}

func TestProcessor_MarkTxWithFakeToken(t *testing.T) {
	//p := Processor{}
	d := dao.NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "select id from anyswap a where direction = 'in' and match_id is null and \n(from_chain = '56' or from_chain ='137' or from_chain ='250' or from_chain ='42161' or from_chain = '43114')\nand (to_chain='56' or to_chain='137' or to_chain='250' or to_chain='42161' or to_chain='43114')"
	id := []*uint64{}
	err := d.DB().Select(&id, stmt)
	if err != nil {
		fmt.Println(err)
	}

}

func TestProcessor_UpdateProfitAndRisk(t *testing.T) {
	d := dao.NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "select * from anyswap a where profit::jsonb @>'[]' and tag is null"
	datas := []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}

	println(len(datas))
	m := Processor{}
	i := 0
	size := len(datas) / 5
	for ; i < len(datas)-2*size; i = i + size {
		println(i)
		go m.UpdateProfitAndRisk(d, datas[i:i+size])
	}
	println(i)
	m.UpdateProfitAndRisk(d, datas[i:])
}

func TestProcessor_UpdateRisk(t *testing.T) {
	d := dao.NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "select * from anyswap a where profit::jsonb @>'[]' and (tag is null or tag = '2')"
	//stmt := "select * from anyswap where id = 4181095"

	datas := []*model.Data{}
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}

	println(len(datas))
	m := Processor{}
	i := 0
	size := len(datas) / 100

	for i = 0; i+2*size < len(datas); i = i + size {
		go m.UpdateRisk(d, datas[i:i+size])
	}
	m.UpdateRisk(d, datas[i:])
}

func TestP(t *testing.T) {
	stmt := "select * from anyswap where direction = 'out' and match_tag =''"
	data := []*model.Data{}
	err := d.DB().Select(&data, stmt)
	if err != nil {
		fmt.Println(err)
	}
	println(len(data))

	//size := len(data) / 200
	//i := 0
	/*
		for ; i < len(data)-2*size; i = i + size {
			go f(data[i : i+size])
		}*/
	f(data)
}

func f(data model.Datas) {
	for _, i := range data {
		stmt := fmt.Sprintf("update anyswap set match_tag = hash where id=%d", i.Id)
		_, err := d.DB().Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("done: ", data[len(data)-1].Id)
}
