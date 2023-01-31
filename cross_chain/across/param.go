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
	"eth": {
		"0x4d9079bb4165aeb4084c526a32695dcfd2f77381",
	},
	"arbitrum": {
		"0xB88690461dDbaB6f04Dfad7df66B7725942FEb9C",
	},
	"polygon": {
		"0x69B5c72837769eF1e7C164Abc6515DcFf217F920",
	},
	"optimism": {
		"0xa420b2d1c0841415A695b81E5B867BCD07Dff8C9",
	},
	"boba": {
		"0xBbc6009fEfFc27ce705322832Cb2068F8C1e0A58",
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
