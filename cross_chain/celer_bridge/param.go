package celer_bridge

import (
	"app/utils"
	"math/big"
)

const (
	//Burn_1 (bytes32 burnId, address token, address account, uint256 amount, address withdrawAccount)
	Burn_1 = "0x75f1bf55bb1de41b63a775dc7d4500f01114ee62b688a6b11d34f4692c1f3d43"
	//Burn_2 (bytes32 burnId, address token, address account, uint256 amount, uint64 toChainId, address toAccount, uint64 nonce)
	Burn_2 = "0x6298d7b58f235730b3b399dc5c282f15dae8b022e5fbbf89cee21fd83c8810a3"
	//Send (bytes32 transferId, address sender, address receiver, address token, uint256 amount, uint64 dstChainId, uint64 nonce, uint32 maxSlippage)
	Send = "0x89d8051e597ab4178a863a5190407b98abfeff406aa8db90c59af76612e58f01"
	//Deposited_1 (bytes32 depositId, address depositor, address token, uint256 amount, uint64 mintChainId, address mintAccount)
	Deposited_1 = "0x15d2eeefbe4963b5b2178f239ddcc730dda55f1c23c22efb79ded0eb854ac789"
	//Deposited (bytes32 depositId, address depositor, address token, uint256 amount, uint64 mintChainId, address mintAccount, uint64 nonce)
	Deposited_2 = "0x28d226819e371600e26624ebc4a9a3947117ee2760209f816c789d3a99bf481b"

	//Burn ==> Mint
	//Send ==> Relay
	//Deposited_1 & Deposited_2 ==> Withdrawn

	//Mint (bytes32 mintId, address token, address account, uint256 amount, uint64 refChainId, bytes32 refId, address depositor)
	Mint = "0x5bc84ecccfced5bb04bfc7f3efcdbe7f5cd21949ef146811b4d1967fe41f777a"
	//Relay (bytes32 transferId, address sender, address receiver, address token, uint256 amount, uint64 srcChainId, bytes32 srcTransferId)
	Relay = "0x79fa08de5149d912dce8e5e8da7a7c17ccdf23dd5d3bfe196802e6eb86347c7c"
	//Withdrawn (bytes32 withdrawId, address receiver, address token, uint256 amount, uint64 refChainId, bytes32 refId, address burnAccount)
	Withdrawn = "0x296a629c5265cb4e5319803d016902eb70a9079b89655fe2b7737821ed88beeb"
)

var CBridgeContracts = map[string][]string{
	"eth": {
		"0x5427FEFA711Eff984124bFBB1AB6fbf5E3DA1820",
		"0xB37D31b2A74029B5951a2778F959282E2D518595",
		"0x7510792A3B1969F9307F3845CE88e39578f2bAE1",
		"0x52E4f244f380f8fA51816c8a10A63105dd4De084",
		"0x16365b45EB269B5B5dACB34B4a15399Ec79b95eB",
	},
	"arbitrum": {
		"0xb3833Ecd19D4Ff964fA7bc3f8aC070ad5e360E56",
		"0x1619DE6B6B20eD217a58d00f37B9d47C7663feca",
		"0xFe31bFc4f7C9b69246a6dc0087D91a91Cb040f76",
		"0xEA4B1b0aa3C110c55f650d28159Ce4AD43a4a58b",
		"0xbdd2739AE69A054895Be33A22b2D2ed71a1DE778",
	},
	"polygon": {
		"0x88DCDC47D2f83a99CF0000FDF667A468bB958a78",
		"0xc1a2D967DfAa6A10f3461bc21864C23C1DD51EeA",
		"0x4C882ec256823eE773B25b414d36F92ef58a7c0C",
		"0xb51541df05DE07be38dcfc4a80c05389A54502BB",
		"0x4d58FDC7d0Ee9b674F49a0ADE11F26C3c9426F7A",
	},
	"optimism": {
		"0x9D39Fc627A6d9d9F8C831c16995b209548cc3401",
		"0xbCfeF6Bb4597e724D720735d32A9249E0640aA11",
		"0x61f85fF2a2f4289Be4bb9B72Fc7010B3142B5f41",
	},
	"fantom": {
		"0x374B8a9f3eC5eB2D97ECA84Ea27aCa45aa1C57EF",
		"0x7D91603E79EA89149BAf73C9038c51669D8F03E9",
		"0x30F7Aa65d04d289cE319e88193A33A8eB1857fb9",
		"0x38D1e20B0039bFBEEf4096be00175227F8939E51",
	},
	"bsc": {
		"0xdd90E5E87A2081Dcf0391920868eBc2FFB81a1aF",
		"0x78bc5Ee9F11d133A08b331C2e18fE81BE0Ed02DC",
		"0x11a0c9270D88C99e221360BCA50c2f6Fda44A980",
		"0x26c76F7FeF00e02a5DD4B5Cc8a0f717eB61e1E4b",
		"0xd443FE6bf23A4C9B78312391A30ff881a097580E",
	},
	"avalanche": {
		"0xef3c714c9425a8F3697A9C969Dc1af30ba82e5d4",
		"0x5427FEFA711Eff984124bFBB1AB6fbf5E3DA1820",
		"0xb51541df05DE07be38dcfc4a80c05389A54502BB",
		"0xb774C6f82d1d5dBD36894762330809e512feD195",
		"0x88DCDC47D2f83a99CF0000FDF667A468bB958a78",
	},
}

func init() {
	for name, chain := range CBridgeContracts {
		CBridgeContracts[name] = utils.StrSliceToLower(chain)
	}
}

type Detail struct {
	TxId        string  `json:"txId,omitempty"`
	Nonce       big.Int `json:"nonce,omitempty"`
	MaxSlippage big.Int `json:"maxSlippage,omitempty"`
}
