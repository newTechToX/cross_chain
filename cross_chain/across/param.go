package across

import (
	"app/utils"
)

const (
	// FundsDeposited (uint256 amount, uint256 originChainId, uint256 destinationChainId, uint64 relayerFeePct, index_topic_1 uint32 depositId, uint32 quoteTimestamp, index_topic_2 address originToken, address recipient, index_topic_3 address depositor)
	FundsDeposited = "0x4a4fc49abd237bfd7f4ac82d6c7a284c69daaea5154430cff04ad7482c6c4254"

	// LogAcrossIn (index_topic_1 bytes32 txhash, index_topic_2 address token, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)
	FilledRelay = "0x56450a30040c51955338a4a9fbafcf94f7ca4b75f4cd83c2f5e29ef77fbe0a3a"
)

// LogAcrossOut (index_topic_1 address token, index_topic_2 address from, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)

var AcrossContracts = map[string][]string{
	"ethereum": {
		"0x4d9079bb4165aeb4084c526a32695dcfd2f77381",
		"0x931a43528779034ac9eb77df799d133557406176",
	},
	"arbitrum": {
		"0xB88690461dDbaB6f04Dfad7df66B7725942FEb9C",
		"0xe1c367e2b576ac421a9f46c9cc624935730c36aa",
	},
	"polygon": {
		"0x69B5c72837769eF1e7C164Abc6515DcFf217F920",
		"0xd3ddacae5afb00f9b9cd36ef0ed7115d7f0b584c",
	},
	"optimism": {
		"0xa420b2d1c0841415A695b81E5B867BCD07Dff8C9",
		"0x59485d57eecc4058f7831f46ee83a7078276b4ae",
	},
	"boba": {
		"0xBbc6009fEfFc27ce705322832Cb2068F8C1e0A58",
		"0x7229405a2f0c550ce35182ee1658302b65672443",
	},
}

func init() {
	for name, chain := range AcrossContracts {
		AcrossContracts[name] = utils.StrSliceToLower(chain)
	}
}

type Detail struct {
	DepositId string `json:"depositId,omitempty"`
	Relayer   string `json:"relayer,omitempty"`
}

var AcrossToken = map[string][]string{
	"eth": {
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		"0x6b175474e89094c44da98b954eedeac495271d0f",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"0x04fa0d235c4abf4bcf4787af4cf447de572ef828",
		"0xba100000625a3754423978a60c9317c58a424e3d",
		"0x42bbfa2e77757c645eeaad1655e0911a7553efbc",
		"0x44108f0223a3c3028f5fe7aec7f9bb2e66bef82f",
		"0xdac17f958d2ee523a2206206994597c13d831ec7",
		"0xc011a73ee8576fb46f5e1c5751ca3b9fe0af2a6f",
		"0x3472a5a71965499acd81997a54bba8d852c6e53d",
	},
	"optimism": {
		"0x4200000000000000000000000000000000000006",
		"0x7f5c764cbc14f9669b88837ca1490cca17c31607",
		"0xda10009cbd5d07dd0cecc66161fc93d7c9000da1",
		"0x68f180fcce6836688e9084f035309e29bf0a2095",
		"0xfe8b128ba8c78aabc59d4c64cee7ff28e9379921",
		"0xe7798f023fc62146e8aa1b36da45fb70855a77ea",
		"0x94b008aa00579c1307b0ef2c499ad98a8ce58e58",
		"0xff733b2a3557a7ed6697007ab5d11b79fdd1b76b",
		"0x8700daec35af8ff88c16bdf0418774cb3d7599b4",
	},
}
