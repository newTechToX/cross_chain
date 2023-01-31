package renbridge

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var renAbi abi.ABI

func init() {
	var err error
	renAbi, err = abi.JSON(strings.NewReader(renAbiStr))
	if err != nil {
		panic(err)
	}
}

func Decode(selector string, hexInput string) ([]interface{}, error) {
	m, err := renAbi.MethodById(common.FromHex(selector))
	if err != nil {
		return nil, err
	}
	return m.Inputs.Unpack(common.FromHex(hexInput))
}
