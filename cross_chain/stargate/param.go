package stargate

import (
	"math/big"
	"strings"
)

const (
	// 1-2 / 1-2-3 / 4
	//SendMsg --- Bridge Contract (uint8 msgType, uint64 nonce)
	SendMsg = "0x8d3ee0df6a4b7e82a7f20a763f1c6826e6176323e655af64f32318827d2112d4"
	//Swap --- Address Contract (uint16 chainId, uint256 dstPoolId, address from, uint256 amountSD, uint256 eqReward, uint256 eqFee, uint256 protocolFee, uint256 lpFee)
	Swap = "0x34660fc8af304464529f48a778e03d03e4d34bcd5f9b6f0cfbf3cd238c642f7f"
	//RedeemRemote --- Address Contract (uint16 chainId, uint256 dstPoolId, address from, uint256 amountLP, uint256 amountSD)
	RedeemRemote = "0xa33f5c0b76f00f6737b1780a8a7f18e19c3fe8fe9ee01a6c1b8ce1eae5ed54f9"

	//SendToChain --- StargateToken Contract (uint16 dstChainId, bytes to, uint256 qty)
	SendToChain = "0x664e26797cde1146ddfcb9a5d3f4de61179f9c11b2698599bb09e686f442172b"

	// 1-2 / 1-3 / 1-4
	//PacketReceived --- Layer Zero Contract (index_topic_1 uint16 srcChainId, bytes srcAddress, index_topic_2 address dstAddress, uint64 nonce, bytes32 payloadHash)
	PacketReceived = "0x2bd2d8a84b748439fd50d79a49502b4eb5faa25b864da6a9ab5c150704be9a4d"
	//SwapRemote --- Address Contract (address to, uint256 amountSD, uint256 protocolFee, uint256 dstFee)
	SwapRemote = "0xfb2b592367452f1c437675bed47f5e1e6c25188c17d7ba01a12eb030bc41ccef"
	//RedeemLocalCallback --- Address Contract (uint16 srcChainId, index_topic_1 bytes srcAddress, index_topic_2 uint256 nonce, uint256 srcPoolId, uint256 dstPoolId, address to, uint256 amountSD, uint256 mintAmountSD)
	RedeemLocalCallback = "0xc7379a02e530fbd0a46ea1ce6fd91987e96535798231a796bdc0e1a688a50873"
	//ReceiveFromChain --- StargateToken Contract (uint16 srcChainId, uint64 nonce, uint256 qty)
	ReceiveFromChain = "0x831bc68226f8d1f734ffcca73602efc4eca13711402ba1d2cc05ee17bb54f631"
)

// PacketReceived 与 SwapRemote / ReceivedFromChain / RedeemLocalCallback 会在同一笔交易中产生
// SendMsg 与 Swap / RedeemRemote(?) 会在同一笔交易中产生
// SendToChain 只会自己产生，目标链上对应 ReceivedFromChain

var StargateContracts = map[string][]string{
	"eth": {
		"0x8731d54E9D02c286767d56ac03e8037C07e01e98",
		"0x150f94B44927F078737562f0fcF3C95c01Cc2376",
		"0xAf5191B0De278C7286d6C7CC6ab6BB8A73bA2Cd6",
		"0x38EA452219524Bb87e18dE1C24D3bB59510BD783",
		"0x8731d54E9D02c286767d56ac03e8037C07e01e98",
		"0x692953e758c3669290cb1677180c64183cEe374e",
		"0x0Faf1d2d3CED330824de3B8200fc8dc6E397850d",
		"0xfa0f307783ac21c39e939acff795e27b650f6e68",
		"0x590d4f8A68583639f215f675F3a259Ed84790580",
		"0xE8F55368C82D38bbbbDb5533e7F56AfC2E978CC2",
		"0x9cef9a0b1bE0D289ac9f4a98ff317c33EAA84eb8",
		"0xdf0770df86a8034b3efef0a1bb3c889b8332ff56",
		"0x296F55F8Fb28E498B858d0BcDA06D955B2Cb3f97",
		"0x101816545f6bd2b1076434b54383a1e633390a2e",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2", // Layer Zero Contract
	},
	"arbitrum": {
		"0x53Bf833A5d6c4ddA888F69c22C88C9f356a41614",
		"0xbf22f0f184bCcbeA268dF387a49fF5238dD23E40",
		"0x6694340fc020c5E6B96567843da2df01b2CE1eb6",
		"0x915A55e36A01285A14f05dE6e81ED9cE89772f8e",
		"0x892785f33CdeE22A30AEF750F285E18c18040c3e",
		"0xB6CfcF89a7B22988bfC96632aC2A9D6daB60d641",
		"0xaa4BF442F024820B2C28Cd0FD72b82c63e66F56C",
		"0xF39B7Be294cB36dE8c510e267B82bb588705d977",
		"0x352d8275AAE3e0c2404d9f68f6cEE084B5bEB3DD",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2",
	},
	"polygon": {
		"0x45A01E4e04F14f7A4a6702c74187c5F6222033cd",
		"0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590",
		"0x1205f31718499dBf1fCa446663B532Ef87481fe1",
		"0x29e38769f23701A2e4A8Ef0492e19dA4604Be62c",
		"0x1c272232Df0bb6225dA87f4dEcD9d37c32f63Eea",
		"0x8736f92646B2542B3e5F3c63590cA7Fe313e283B",
		"0x9d1B1669c73b033DFe47ae5a0164Ab96df25B944",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2",
	},
	"optimism": {
		"0xB0D502E938ed5f4df2E681fE6E419ff29631d62b",
		"0xB49c4e680174E331CB0A7fF3Ab58afC9738d5F8b",
		"0x296F55F8Fb28E498B858d0BcDA06D955B2Cb3f97",
		"0xd22363e3762cA7339569F3d33EADe20127D5F98C",
		"0xDecC0c09c3B5f6e92EF4184125D5648a66E35298",
		"0x165137624F1f692e69659f944BF69DE02874ee27",
		"0x368605D9C6243A80903b9e326f1Cddde088B8924",
		"0x2F8bC9081c7FCFeC25b9f41a50d97EaA592058ae",
		"0x3533F5e279bDBf550272a199a223dA798D9eff78",
		"0x5421FA1A48f9FF81e4580557E86C7C0D24C18036",
		"0x701a95707A0290AC8B90b3719e8EE5b210360883",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2",
	},
	"fantom": {
		"0xAf5191B0De278C7286d6C7CC6ab6BB8A73bA2Cd6",
		"0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590",
		"0x12edeA9cd262006cC3C4E77c90d2CD2DD4b1eb97",
		"0x45A01E4e04F14f7A4a6702c74187c5F6222033cd",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2",
	},
	"bsc": {
		"0x6694340fc020c5E6B96567843da2df01b2CE1eb6",
		"0x9aa83081aa06af7208dcc7a4cb72c94d057d2cda",
		"0x98a5737749490856b401DB5Dc27F522fC314A4e1",
		"0x4e145a589e4c03cBe3d28520e4BF3089834289Df",
		"0x7BfD7f2498C4796f10b6C611D9db393D3052510C",
		"0x6694340fc020c5E6B96567843da2df01b2CE1eb6",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2",
		"0xB0D502E938ed5f4df2E681fE6E419ff29631d62b",
	},
	"avalanche": {
		"0x45A01E4e04F14f7A4a6702c74187c5F6222033cd",
		"0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590",
		"0x1205f31718499dBf1fCa446663B532Ef87481fe1",
		"0x29e38769f23701A2e4A8Ef0492e19dA4604Be62c",
		"0x1c272232Df0bb6225dA87f4dEcD9d37c32f63Eea",
		"0x8736f92646B2542B3e5F3c63590cA7Fe313e283B",
		"0x4d73adb72bc3dd368966edd0f0b2148401a178e2",
		"0x9d1B1669c73b033DFe47ae5a0164Ab96df25B944",
	},
}

func init() {
	for name, chain := range StargateContracts {
		StargateContracts[name] = ToLower(chain)
	}
}

func ToLower(s []string) []string {
	ret := make([]string, 0)
	for _, r := range s {
		ret = append(ret, strings.ToLower(r))
	}
	return ret
}

type Detail struct {
	Nonce *big.Int `json:"nonce,omitempty"`
}
