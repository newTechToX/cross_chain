package utils

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

func ParseStrToUint64(nStr string) uint64 {
	if strings.HasPrefix(nStr, "0x") || strings.HasPrefix(nStr, "0X") {
		if len(nStr) == 2 {
			return 0
		}
		bigVal, ok := new(big.Int).SetString(nStr[2:], 16)
		if !ok || bigVal == nil {
			log.Error("invalid hex", "number", nStr)
			return 0
		}
		return bigVal.Uint64()
	}
	val, err := strconv.ParseUint(nStr, 10, 64)
	if err == nil {
		return val
	}
	bigVal, ok := new(big.Int).SetString(nStr, 16)
	if ok && bigVal != nil {
		return bigVal.Uint64()
	}
	if len(nStr) == 0 {
		return 0
	}
	log.Error("invalid number", "number", nStr)
	return 0
}
