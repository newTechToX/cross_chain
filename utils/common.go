package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

const (
	EtherScanMaxResult = 10000
)

func PrintPretty(data interface{}) {
	res, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(res))
}

func StrSliceToLower(s []string) []string {
	ret := make([]string, 0)
	for _, r := range s {
		ret = append(ret, strings.ToLower(r))
	}
	return ret
}

func HexSum(hexes ...string) *big.Int {
	ret := new(big.Int).SetUint64(0)
	for _, hex := range hexes {
		t, _ := new(big.Int).SetString(strings.TrimPrefix(hex, "0x"), 16)
		ret.Add(ret, t)
	}
	return ret
}

func Contains[T constraints.Ordered](target T, slice []T) bool {
	for _, e := range slice {
		if target == e {
			return true
		}
	}
	return false
}

//
func ParseDateTime(ts string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", ts)
}

func IsTargetCall(input string, selectors []string) bool {
	if len(selectors) == 0 {
		return true
	}
	for _, s := range selectors {
		if len(s) != 10 {
			continue
		}
		if strings.HasPrefix(input, s) {
			return true
		}
	}
	return false
}

func DeleteSliceElementByValue[T constraints.Ordered](s []T, e T) []T {
	i := 0
	for _, v := range s {
		if v != e {
			s[i] = v
			i++
		}
	}
	return s[:i]
}
