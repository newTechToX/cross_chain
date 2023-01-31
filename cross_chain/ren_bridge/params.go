package renbridge

import (
	"app/utils"
	"math/big"
)

const (
	// burn(bytes _to, uint256 _amount)
	Burn = "0x38463cff"

	// mint(bytes32 _pHash, uint256 _amountUnderlying, bytes32 _nHash, bytes _sig)
	Mint = "0x159ab14d"
)

var contracts = map[string]map[string]struct {
	Token   string
	ChainId *big.Int
}{
	"eth": {
		"0xe4b679400f0f267212d5d812b95f58c83243ee71": {"0xeb4c2781e4eba804ce9a9803c67d0893436bb27d", utils.GetChainId("btc")},
	},
	"bsc": {
		"0x95de7b32e24b62c44a4c44521eff4493f1d1fe13": {"0xfce146bf3146100cfe5db4129cf6c82b0ef4ad8c", utils.GetChainId("btc")},
	},
	"polygon": {
		"0x05cadbf3128bcb7f2b89f3dd55e5b0a036a49e20": {"0xdbf31df14b66535af65aac99c32e9ea844e14501", utils.GetChainId("btc")},
	},
	"avalanche": {
		"0x05cadbf3128bcb7f2b89f3dd55e5b0a036a49e20": {"0xdbf31df14b66535af65aac99c32e9ea844e14501", utils.GetChainId("btc")},
	},
	"arbitrum": {
		"0x05cadbf3128bcb7f2b89f3dd55e5b0a036a49e20": {"0xdbf31df14b66535af65aac99c32e9ea844e14501", utils.GetChainId("btc")},
	},
	"optimism": {
		"0xb538901719936e628a9b9af64a5a4dbc273305cd": {"0x85f6583762bc76d775eab9a7456db344f12409f7", utils.GetChainId("btc")},
	},
	"fantom": {
		"0x05cadbf3128bcb7f2b89f3dd55e5b0a036a49e20": {"0xdbf31df14b66535af65aac99c32e9ea844e14501", utils.GetChainId("btc")},
	},
	"moonbeam": {
		"0xb538901719936e628a9b9af64a5a4dbc273305cd": {"0x85f6583762bc76d775eab9a7456db344f12409f7", utils.GetChainId("btc")},
	},
}
