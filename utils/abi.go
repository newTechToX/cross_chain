package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func DecodeInput(myAbi abi.ABI, selector string, input string) ([]interface{}, error) {
	m, err := myAbi.MethodById(common.FromHex(selector))
	if err != nil {
		return nil, err
	}
	return m.Inputs.Unpack(common.FromHex(input))
}
