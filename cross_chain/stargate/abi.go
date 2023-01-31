package stargate

const (
	stargateAbiStr = `[
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "uint16",
					"name": "srcChainId",
					"type": "uint16"
				},
				{
					"indexed": false,
					"internalType": "bytes",
					"name": "srcAddress",
					"type": "bytes"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "dstAddress",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint64",
					"name": "nonce",
					"type": "uint64"
				},
				{
					"indexed": false,
					"internalType": "bytes32",
					"name": "payloadHash",
					"type": "bytes32"
				}
			],
			"name": "PacketReceived",
			"type": "event"
		}
	]`
)
