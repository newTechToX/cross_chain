package wormhole

const (
	// transferTokens(address token,uint256 amount,uint16 recipientChain,bytes32 recipient,uint256 arbiterFee,uint32 nonce)
	TransferTokens = "0x0f5287b0"

	// wrapAndTransferETH(uint16 recipientChain,bytes32 recipient,uint256 arbiterFee,uint32 nonce)
	WrapAndTransferETH = "0x9981509f"

	// function transferTokensWithPayload( address token, uint256 amount, uint16 recipientChain, bytes32 recipient, uint32 nonce, bytes memory payload)
	TransferTokensWithPayload = "0xc5a5ebda"

	// function wrapAndTransferETHWithPayload( uint16 recipientChain, bytes32 recipient, uint32 nonce, bytes memory payload )
	WrapAndTransferETHWithPayload = "0xbee9cdfc"

	// completeTransfer(bytes encodedVm)
	CompleteTransfer = "0xc6878519"

	// completeTransferAndUnwrapETH(bytes encodedVm)
	CompleteTransferAndUnwrapETH = "0xff200cde"

	// completeTransferWithPayload(bytes memory encodedVm)
	CompleteTransferWithPayload = "0xc3f511c1"

	// completeTransferAndUnwrapETHWithPayload(bytes memory encodedVm)
	CompleteTransferAndUnwrapETHWithPayload = "0x1c8475e4"
)

var contracts = map[string]map[string]struct{}{
	"eth": {
		"0x3ee18b2214aff97000d974cf647e7c347e8fa585": {},
	},
	"bsc": {
		"0xb6f6d86a8f9879a9c87f643768d9efc38c1da6e7": {},
	},
	"polygon": {
		"0x5a58505a96d1dbf8df91cb21b54419fc36e93fde": {},
	},
	"avalanche": {
		"0x0e082f06ff657d94310cb8ce8b0d9a04541d8052": {},
	},
	"fantom": {
		"0x7c9fc5741288cdfdd83ceb07f3ea7e22618d79d2": {},
	},
}

type OutDetail struct {
	Nonce uint32
}
