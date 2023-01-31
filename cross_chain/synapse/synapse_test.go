package synapse

import (
	"testing"
)

func TestSynapse_Contracts(t *testing.T) {
	a := &Synapse{}
	a.Contracts("eth")
	for addr, _ := range contracts["eth"] {
		println(addr)
	}
}
