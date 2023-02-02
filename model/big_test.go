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

func TestBigFloat_SetString(t *testing.T) {
	s := "-0.567892938473455229384759856"
	e := new(BigFloat).SetString(s)
	ee := e.ConvertToBigInt()
	println(ee.String())
}
