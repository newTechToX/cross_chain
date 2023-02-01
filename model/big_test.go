package model

import (
	"fmt"
	"math/big"
	"testing"
)

func TestBigFloat_Scan(t *testing.T) {

	b, _ := new(big.Float).SetString("111")
	//bb, _ := b.MarshalText()
	//a.SetBigFloat(b)
	//fmt.Println(b.String())
	fmt.Println(((*BigFloat)(b).String()))
}
