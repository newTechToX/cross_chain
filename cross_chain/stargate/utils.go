package stargate

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var stargateAbi abi.ABI

func init() {
	var err error
	stargateAbi, err = abi.JSON(strings.NewReader(stargateAbiStr))
	if err != nil {
		panic(err)
	}
}

func DecodePacketReceivedData(data string) ([]interface{}, error) {
	return stargateAbi.Unpack("PacketReceived", common.FromHex(data))
}
