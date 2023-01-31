package synapse

import (
	"strings"
)

const (
	//TokenRedeem (index_topic_1 address to, uint256 chainId, address token, uint256 amount)
	TokenRedeem = "0xdc5bad4651c5fbe9977a696aadc65996c468cde1448dd468ec0d83bf61c4b57c"
	//TokenRedeemAndSwap (index_topic_1 address to, uint256 chainId, address token, uint256 amount, uint8 tokenIndexFrom, uint8 tokenIndexTo, uint256 minDy, uint256 deadline)
	TokenRedeemAndSwap = "0x91f25e9be0134ec851830e0e76dc71e06f9dade75a9b84e9524071dbbc319425"
	//TokenRedeemAndRemove
	TokenRedeemAndRemove = "0x9a7024cde1920aa50cdde09ca396229e8c4d530d5cfdc6233590def70a94408c"
	//TokenDeposit (index_topic_1 address to, uint256 chainId, address token, uint256 amount)
	TokenDeposit = "0xda5273705dbef4bf1b902a131c2eac086b7e1476a8ab0cb4da08af1fe1bd8e3b"
	//TokenDepositAndSwap (index_topic_1 address to, uint256 chainId, address token, uint256 amount, uint8 tokenIndexFrom, uint8 tokenIndexTo, uint256 minDy, uint256 deadline)
	TokenDepositAndSwap = "0x79c15604b92ef54d3f61f0c40caab8857927ca3d5092367163b4562c1699eb5f"

	//TokenMint (index_topic_1 address to, address token, uint256 amount, uint256 fee, index_topic_2 bytes32 kappa)
	TokenMint = "0xbf14b9fde87f6e1c29a7e0787ad1d0d64b4648d8ae63da21524d9fd0f283dd38"
	//TokenMintAndSwap (index_topic_1 address to, address token, uint256 amount, uint256 fee, uint8 tokenIndexFrom, uint8 tokenIndexTo, uint256 minDy, uint256 deadline, bool swapSuccess, index_topic_2 bytes32 kappa)
	TokenMintAndSwap = "0x4f56ec39e98539920503fd54ee56ae0cbebe9eb15aa778f18de67701eeae7c65"
	//TokenWithdraw (index_topic_1 address to, address token, uint256 amount, uint256 fee, index_topic_2 bytes32 kappa)
	TokenWithdraw = "0x8b0afdc777af6946e53045a4a75212769075d30455a212ac51c9b16f9c5c9b26"
	//TokenWithdrawAndRemove (index_topic_1 address to, address token, uint256 amount, uint256 fee, uint8 swapTokenIndex, uint256 swapMinAmount, uint256 swapDeadline, bool swapSuccess, index_topic_2 bytes32 kappa)
	TokenWithdrawAndRemove = "0xc1a608d0f8122d014d03cc915a91d98cef4ebaf31ea3552320430cba05211b6d"
)

const (
	deposit                 = "0x90d25074"
	depositAndSwap          = "0xa2a2af0b"
	depositETH              = "0xce0b63ce"
	depositETHAndSwap       = "0xa9fd8c57"
	redeem                  = "0xf3f094a1"
	redeemv2                = "0x49b7cf84"
	redeemAndSwap           = "0x839ed90a"
	redeemAndRemove         = "0x36e712ed"
	swapAndRedeem           = "0x4a517a55"
	swapAndRedeemAndRemove  = "0x174dc952"
	swapAndRedeemAndSwap    = "0x9f330727"
	swapETHAndRedeem        = "0x65749c9d"
	swapETHAndRedeemAndSwap = "0x3d5da164"
	zapAndDeposit           = "0x394023d9"
	zapAndDepositAndSwap    = "0x8054ae3b"
)

var contracts = map[string]map[string]struct{}{
	"eth": {
		"0x2796317b0fF8538F253012862c06787Adfb8cEb6": {},
		"0x6571d6be3d8460CF5F7d6711Cd9961860029D85F": {},
	},
	"arbitrum": {
		"0x37f9ae2e0ea6742b9cad5abcfb6bbc3475b3862b": {},
		"0x6f4e8eba4d337f874ab57478acc2cb5bacdc19c9": {},
	},
	"polygon": {
		"0x1c6ae197ff4bf7ba96c66c5fd64cb22450af9cc8": {},
		"0x8f5bbb2bb8c2ee94639e55d5f41de9b4839c1280": {},
	},
	"optimism": {
		"0x470f9522ff620ee45df86c58e54e6a645fe3b4a7": {},
		"0xaf41a65f786339e7911f4acdad6bd49426f2dc6b": {},
	},
	"fantom": {
		"0xb003e75f7e0b5365e814302192e99b4ee08c0ded": {},
		"0xaf41a65f786339e7911f4acdad6bd49426f2dc6b": {},
	},
	"bsc": {
		"0xd123f70ae324d34a9e76b67a27bf77593ba8749f": {},
		"0x749f37df06a99d6a8e065dd065f8cf947ca23697": {},
	},
	"avalanche": {
		"0x0ef812f4c68dc84c22a4821ef30ba2ffab9c2f3a": {},
		"0xc05e61d0e7a63d27546389b7ad62fdff5a91aace": {},
	},
}

func init() {
	for _, addrs := range contracts {
		for addr, _ := range addrs {
			addr = strings.ToLower(addr)
		}
	}
}

type Detail struct {
	Kappa string `json:"kappa,omitempty"`
}
