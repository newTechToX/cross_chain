package anyswap

import "app/utils"

const (
	// LogAnySwapOut (index_topic_1 address token, index_topic_2 address from, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)
	LogAnySwapOut = "0x97116cf6cd4f6412bb47914d6db18da9e16ab2142f543b86e207c24fbd16b23a"

	// LogAnySwapIn (index_topic_1 bytes32 txhash, index_topic_2 address token, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)
	LogAnySwapIn = "0xaac9ce45fe3adf5143598c4f18a369591a20a3384aedaf1b525d29127e1fcd55"

	// underlying()
	Underlying = "0x6f307dc3"
)

// LogAnySwapOut (index_topic_1 address token, index_topic_2 address from, index_topic_3 address to, uint256 amount, uint256 fromChainID, uint256 toChainID)

var AnyswapContracts = map[string][]string{
	"eth": {
		"0x6b7a87899490ece95443e979ca9485cbe7e71522",
		"0xba8da9dcf11b50b03fd5284f164ef5cdef910705",
		"0x765277eebeca2e31912c9946eae1021199b39c61",
		"0x7782046601e7b9b05ca55a3899780ce6ee6b8b2b",
		"0xe95fd76cf16008c12ff3b3a937cb16cd9cc20284",
		"0xf0457c4c99732b716e40d456acb3fc83c699b8ba",
	},
	"bsc": {
		"0xd1c5966f9f5ee6881ff6b261bbeda45972b1b5f3",
		"0xabd380327fe66724ffda91a87c772fb8d00be488",
		"0x56a6c850cebe23f0c7891a004bef57265cda4d13",
		"0x58892974758a4013377a45fad698d2ff1f08d98e",
		"0x92c079d3155c2722dbf7e65017a5baf9cd15561c",
		"0xd1a891e6eccb7471ebd6bc352f57150d4365db21",
		"0xe1d592c3322f1f714ca11f05b6bc0efef1907859",
		"0xf9736ec3926703e85c843fc972bd89a7f8e827c0",
	},
	"polygon": {
		"0x4f3aff3a747fcade12598081e80c6605a8be192f",
		"0x2ef4a574b72e1f555185afa8a09c6d1a8ac4025c",
		"0x0b23341fa1da0171f52aa8ef85f3946b44d35ac0",
		"0x1ccca1ce62c62f7be95d4a67722a8fdbed6eecb4",
		"0x6ff0609046a38d76bd40c5863b4d1a2dce687f73",
		"0x72c290f3f13664b024ee611983aa2d5621ebe917",
		"0x84cebca6bd17fe11f7864f7003a1a30f2852b1dc",
		"0xafaace7138ab3c2bcb2db4264f8312e1bbb80653",
		"0xd50380e953603b37a74dc67c92fc5e19e0b65469",
	},
	"fantom": {
		"0x1ccca1ce62c62f7be95d4a67722a8fdbed6eecb4",
		"0x0b23341fa1da0171f52aa8ef85f3946b44d35ac0",
		"0x24e2a6f08e3cc2baba93bd9b89e19167a37d6694",
		"0x85fd5f8dbd0c9ef1806e6c7d4b787d438621c1dc",
		"0xb576c9403f39829565bd6051695e2ac7ecf850e2",
		"0xf3ce95ec61114a4b1bfc615c16e6726015913ccc",
		"0xf98f70c265093a3b3adbef84ddc29eace900685b",
	},
	"arbitrum": {
		"0x0cae51e1032e8461f4806e26332c030e34de3adb",
		"0xc931f61b1534eb21d8c11b24f3f5ab2471d4ab50",
		"0x650af55d5877f289837c30b94af91538a7504b76",
		"0x2bf9b864cdc97b08b6d79ad4663e71b8ab65c45c",
		"0x39fde572a18448f8139b7788099f0a0740f51205",
		"0xa71353bb71dda105d383b02fc2dd172c4d39ef8b",
		"0xcb9f441ffae898e7a2f32143fd79ac899517a9dc",
	},
	"avalanche": {
		"0xb0731d50c681c45856bfc3f7539d5f61d4be81d8",
		"0x833f307ac507d47309fd8cdd1f835bef8d702a93",
		"0x05f024c6f5a94990d32191d6f36211e3ee33504e",
		"0x34324e1598bf02ccd3dea93f4e332b5507097473",
		"0x9b17baadf0f21f03e35249e0e59723f34994f806",
		"0xe5cf1558a1470cb5c166c2e8651ed0f3c5fb8f42",
	},
	"optimism": {
		"0x80a16016cc4a2e6a2caca8a4a498b1699ff0f844",
		"0xdc42728b0ea910349ed3c6e1c9dc06b5fb591f98",
	},
}

var AnyTokens = map[string]map[string]map[string]string{
	"USDC": {
		"eth": {
			"underlyingToken": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			"anyToken_1":      "0x7EA2be2df7BA6E54B1A9C70676f668455E329d29",
			"anyToken_2":      "0xeA928a8d09E11c66e074fBf2f6804E19821F438D",
			"anyToken_3":      "0x2cb1712fa24aBc7Ce787b8853235C86e38ACca44",
		},
		"bsc": {
			"underlyingToken": "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d",
			"anyToken":        "0x8965349fb649A33a30cbFDa057D8eC2C48AbE2A2",
			"anyToken_2":      "0xab6290bBd5C2d26881E8A7a10bC98552B9082E7f",
		},
		"avanlanche": {
			"underlyingToken": "0xA7D7079b0FEaD91F3e65f86E8915Cb59c1a4C664",
			"anyToken":        "0xcc9b1F919282c255eB9AD2C0757E8036165e0cAd",
		},
		"polygon": {
			"underlyingToken": "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
			"anyToken":        "0xd69b31c3225728CC57ddaf9be532a4ee1620Be51",
		},
		"fantom": {
			"underlyingToken": "0x04068DA6C83AFCFA0e13ba15A6696662335D5B75",
			"anyToken":        "0x95bf7E307BC1ab0BA38ae10fc27084bC36FcD605",
		},
		"arbitrum": {
			"underlyingToken": "0x04068DA6C83AFCFA0e13ba15A6696662335D5B75",
			"anyToken":        "0x3405A1bd46B85c5C029483FbECf2F3E611026e45",
		},
		"optimism": {
			"underlyingToken": "0x7F5c764cBc14f9669B88837ca1490cCa17c31607",
			"anyToken":        "0xf390830DF829cf22c53c8840554B98eafC5dCBc2",
		},
	},
}

func init() {
	for name, chain := range AnyswapContracts {
		AnyswapContracts[name] = utils.StrSliceToLower(chain)
	}
}

type Detail struct {
	SrcTxHash string `json:"src_tx_hash,omitempty"`
}

var AnyDeployer = map[string][]string{
	"eth": {
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
		"0xfd63bf84471bc55dd9a83fdfa293ccbd27e1f4c8",
	},
	"bsc": {
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
		"0x790deeb2929201067a460974612c996d2a25183d",
		"0xcf799767d366d789e8b446981c2d578e241fa25c",
	},
	"arbitrum": {
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
		"0xcf799767d366d789e8b446981c2d578e241fa25c",
	},
	"polygon": {
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
	},
	"fantom": {
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
	},
	"avalanche": {
		"0xfd63bf84471bc55dd9a83fdfa293ccbd27e1f4c8",
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
	},
	"optimism": {
		"0xfa9da51631268a30ec3ddd1ccbf46c65fad99251",
	},
}
