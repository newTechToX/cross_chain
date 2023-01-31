package utils

import "errors"

var (
	ErrTooManyRecords     = errors.New("collector: too many records")
	ErrInvalidKey         = errors.New("etherscan: invalid api key")
	ErrEtherscanRateLimit = errors.New("etherscan: rate limit")
)
