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
		"0xdc9524a8774dc2956bdb8b55fdf91938757f3185",
	}
	info, _ := a.QueryAml("eth", address)
	fmt.Println(info["0xdc9524a8774dc2956bdb8b55fdf91938757f3185"][0])
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
