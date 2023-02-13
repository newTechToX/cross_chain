package aml

import (
	"fmt"
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
