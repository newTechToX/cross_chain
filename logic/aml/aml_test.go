package aml

import (
	"app/dao"
	"app/model"
	"fmt"
	"strings"
	"testing"
)

func TestAML_NewAML(t *testing.T) {
	a := NewAML("../txt_config.yaml")
	fmt.Println(a.key)
}

func TestAML_QueryAml(t *testing.T) {
	a := NewAML("../txt_config.yaml")
	address := []string{
		"0x4023ef3aaa0669faaf3a712626f4d8ccc3eaf2e5",
		"0x6ca6568374966713738028c3aed52855ea5e61d3",
		"0xb1d6bc439f5d3bfbd828da3d0848b0f3658c9dc6",
	}
	info, _ := a.QueryAml("eth", address)
	fmt.Println(info["0x8683e604cdf911cd72652a04bf9d571697a86a60"][0])
	fmt.Println(info["0x0b15ddf19d47e6a86a56148fb4afffc6929bcb89"][0])
}

func TestNewAML(t *testing.T) {
	a := []*EntityInfo{
		{
			Entity:     "MULSER",
			EntityType: "sdf",
			Confidence: 9,
		},
	}
	b := []*PropertyInfo{{
		AddressProperty: "swer",
		Confidence:      2},
	}
	c := &Labels{
		EntityInfo:   a,
		PropertyInfo: b,
	}
	ss := c.name()
	println(ss)
}

func TestAML_QueryAml2(t *testing.T) {
	a := NewAML("../txt_config.yaml")
	stmt := "select chain, token from anyswap a where isfaketoken != 1 group by chain, token"
	var d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	var datas = model.TokenChains{}
	if err := d.DB().Select(&datas, stmt); err != nil {
		fmt.Println(err)
	}
	println(len(datas))
	for _, data := range datas {
		info, _ := a.QueryAml(data.Chain, []string{data.Token})
		if v, ok := info[data.Token]; ok {
			for _, dd := range v {
				if strings.Contains(strings.ToLower(dd.Name), "multichain") ||
					strings.Contains(strings.ToLower(dd.Name), "any") {
					stmt = "update anyswap set isfaketoken = 0 where chain = $1 and token = $2"
					if _, err := d.DB().Exec(stmt, dd.Chain, dd.Address); err != nil {
						fmt.Println(err)
					}
				} else {
					println(dd.Chain, dd.Address)
				}
				println(strings.ToLower(dd.Name), dd.Chain, dd.Address)
				for _, e := range dd.Labels {
					println(e)
				}
			}
		}
	}
}
