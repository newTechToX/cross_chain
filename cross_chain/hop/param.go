package hop

import (
	"math/big"
	"strings"
)

const (
	//TransferSentToL2 (index_topic_1 uint256 chainId, index_topic_2 address recipient, uint256 amount, uint256 amountOutMin, uint256 deadline, index_topic_3 address relayer, uint256 relayerFee)
	TransferSentToL2 = "0x0a0607688c86ec1775abcdbab7b33a3a35a6c9cde677c9be880150c231cc6b0b"
	//TransferSent (index_topic_1 bytes32 transferId, index_topic_2 uint256 chainId, index_topic_3 address recipient, uint256 amount, bytes32 transferNonce, uint256 bonderFee, uint256 index, uint256 amountOutMin, uint256 deadline)
	TransferSent = "0xe35dddd4ea75d7e9b3fe93af4f4e40e778c3da4074c9d93e7c6536f1e803c1eb"

	//TransferFromL1Completed (index_topic_1 address recipient, uint256 amount, uint256 amountOutMin, uint256 deadline, index_topic_2 address relayer, uint256 relayerFee)
	TransferFromL1Completed = "0x320958176930804eb66c2343c7343fc0367dc16249590c0f195783bee199d094"
	//WithdrawalBonded (index_topic_1 bytes32 transferId, uint256 amount)
	WithdrawalBonded = "0x0c3d250c7831051e78aa6a56679e590374c7c424415ffe4aa474491def2fe705"
)

var HopContracts = map[string]map[string]map[string]string{

	"USDC": {
		"eth": {
			"canonicalToken": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			"bridge":         "0x3666f603Cc164936C1b87e207F36BEBa4AC5f18a",
		},
		"polygon": {
			"canonicalToken": "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
			"bridge":         "0x25D8039bB044dC227f741a9e381CA4cEAE2E6aE8",
			"ammWrapper":     "0x76b22b8C1079A44F1211D867D68b1eda76a635A7",
		},
		"optimism": {
			"canonicalToken": "0x7F5c764cBc14f9669B88837ca1490cCa17c31607",
			"bridge":         "0xa81D244A1814468C734E5b4101F7b9c0c577a8fC",
			"ammWrapper":     "0x2ad09850b0CA4c7c1B33f5AcD6cBAbCaB5d6e796",
		},
		"arbitrum": {
			"canonicalToken": "0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8",
			"bridge":         "0x0e0E3d2C5c292161999474247956EF542caBF8dd",
			"ammWrapper":     "0xe22D2beDb3Eca35E6397e0C6D62857094aA26F52",
		},
	},

	"USDT": {
		"eth": {
			"canonicalToken": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
			"bridge":         "0x3E4a3a4796d16c0Cd582C382691998f7c06420B6",
		},
		"polygon": {
			"canonicalToken": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
			"bridge":         "0x6c9a1ACF73bd85463A46B0AFc076FBdf602b690B",
			"ammWrapper":     "0x8741Ba6225A6BF91f9D73531A98A89807857a2B3",
		},
		"optimism": {
			"canonicalToken": "0x94b008aA00579c1307B0EF2c499aD98a8ce58e58",
			"bridge":         "0x46ae9BaB8CEA96610807a275EBD36f8e916b5C61",
			"ammWrapper":     "0x7D269D3E0d61A05a0bA976b7DBF8805bF844AF3F",
		},
		"arbitrum": {
			"canonicalToken": "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
			"bridge":         "0x72209Fe68386b37A40d6bCA04f78356fd342491f",
			"ammWrapper":     "0xCB0a4177E0A60247C0ad18Be87f8eDfF6DD30283",
		},
	},

	"MATIC": {
		"eth": {
			"canonicalToken": "0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0",
			"bridge":         "0x22B1Cbb8D98a01a3B71D034BB899775A76Eb1cc2",
		},
		"polygon": {
			"canonicalToken": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
			"bridge":         "0x553bC791D746767166fA3888432038193cEED5E2",
			"ammWrapper":     "0x884d1Aa15F9957E1aEAA86a82a72e49Bc2bfCbe3",
		},
	},

	"DAI": {
		"eth": {
			"canonicalToken": "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			"bridge":         "0x3d4Cc8A61c7528Fd86C55cfe061a78dCBA48EDd1",
		},
		"polygon": {
			"canonicalToken": "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
			"bridge":         "0xEcf268Be00308980B5b3fcd0975D47C4C8e1382a",
			"ammWrapper":     "0x28529fec439cfF6d7D1D5917e956dEE62Cd3BE5c",
		},
		"optimism": {
			"canonicalToken": "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1",
			"bridge":         "0x7191061D5d4C60f598214cC6913502184BAddf18",
			"ammWrapper":     "0xb3C68a491608952Cb1257FC9909a537a0173b63B",
		},
		"arbitrum": {
			"canonicalToken": "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1",
			"bridge":         "0x7aC115536FE3A185100B2c4DE4cb328bf3A58Ba6",
			"ammWrapper":     "0xe7F40BF16AB09f4a6906Ac2CAA4094aD2dA48Cc2",
		},
	},

	"eth": {
		"eth": {
			"canonicalToken": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			"bridge":         "0xb8901acB165ed027E32754E0FFe830802919727f",
		},
		"polygon": {
			"canonicalToken": "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619",
			"bridge":         "0xb98454270065A31D71Bf635F6F7Ee6A518dFb849",
			"ammWrapper":     "0xc315239cFb05F1E130E7E28E603CEa4C014c57f0",
		},
		"optimism": {
			"canonicalToken": "0x4200000000000000000000000000000000000006",
			"bridge":         "0x83f6244Bd87662118d96D9a6D44f09dffF14b30E",
			"ammWrapper":     "0x86cA30bEF97fB651b8d866D45503684b90cb3312",
		},
		"arbitrum": {
			"canonicalToken": "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
			"bridge":         "0x3749C4f034022c39ecafFaBA182555d4508caCCC",
			"ammWrapper":     "0x33ceb27b39d2Bb7D2e61F7564d3Df29344020417",
		},
	},

	"HOP": {
		"eth": {
			"canonicalToken": "0xc5102fE9359FD9a28f877a67E36B0F050d81a3CC",
			"bridge":         "0x914f986a44AcB623A277d6Bd17368171FCbe4273",
		},
		"polygon": {
			"canonicalToken": "0xc5102fE9359FD9a28f877a67E36B0F050d81a3CC",
			"bridge":         "0x58c61AeE5eD3D748a1467085ED2650B697A66234",
			// "ammWrapper":     "0x0000000000000000000000000000000000000000",
		},
		"optimism": {
			"canonicalToken": "0xc5102fE9359FD9a28f877a67E36B0F050d81a3CC",
			"bridge":         "0x03D7f750777eC48d39D080b020D83Eb2CB4e3547",
			// "ammWrapper":     "0x0000000000000000000000000000000000000000",
		},
		"arbitrum": {
			"canonicalToken": "0xc5102fE9359FD9a28f877a67E36B0F050d81a3CC",
			"bridge":         "0x25FB92E505F752F730cAD0Bd4fa17ecE4A384266",
			// "ammWrapper":     "0x0000000000000000000000000000000000000000",
		},
	},

	"SNX": {
		"eth": {
			"canonicalToken": "0xc011a73ee8576fb46f5e1c5751ca3b9fe0af2a6f",
			"bridge":         "0x893246FACF345c99e4235E5A7bbEE7404c988b96",
		},
		"optimism": {
			"canonicalToken": "0x8700dAec35aF8Ff88c16BdF0418774CB3D7599B4",
			"bridge":         "0x16284c7323c35F4960540583998C98B1CfC581a7",
			"ammWrapper":     "0xf11EBB94EC986EA891Aec29cfF151345C83b33Ec",
		},
	},
}

var hopContracts map[string][]string
var hopToken map[string]map[string]string

func init() {
	hopContracts = make(map[string][]string)
	hopToken = make(map[string]map[string]string)

	for _, token := range HopContracts {
		for chainName, chain := range token {
			if _, ok := hopContracts[chainName]; !ok {
				hopContracts[chainName] = []string{}
				hopToken[chainName] = make(map[string]string)
			}

			if chain["bridge"] != "" {
				hopContracts[chainName] = append(hopContracts[chainName], strings.ToLower(chain["bridge"]))
				hopToken[chainName][strings.ToLower(chain["bridge"])] = strings.ToLower(chain["canonicalToken"])

			}

			if chain["ammWrapper"] != "" {
				hopContracts[chainName] = append(hopContracts[chainName], strings.ToLower(chain["ammWrapper"]))
				hopToken[chainName][strings.ToLower(chain["ammWrapper"])] = strings.ToLower(chain["canonicalToken"])
			}
		}
	}
}

type Detail struct {
	DDL        big.Int `json:"ddl,omitempty"`
	TransferID string  `json:"transferID,omitempty"`
	Relayer    string  `json:"relayer,omitempty"`
	MinDy      big.Int `json:"minDy,omitempty"`
}
