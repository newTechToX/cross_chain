package synapse

import (
	"encoding/binary"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

/*var chainId = map[uint64]uint64{
	1:  utils.SolanaChainId,
	2:  1,
	3:  utils.TerraColumbus,
	4:  56,
	5:  137,
	6:  43114,
	7:  4262,
	9:  1313161554,
	10: 250,
	11: 686,
	12: 787,
	13: 8217,
	14: 42220,
	15: utils.NearChainId,
	16: 1284,
	18: utils.TerraPhoenix,
	22: utils.Aptos,
}*/

var synapseAbiETH abi.ABI
var synapseAbiOther abi.ABI

func init() {
	var err error
	synapseAbiETH, err = abi.JSON(strings.NewReader(synapseAbiStrETH))
	synapseAbiOther, err = abi.JSON(strings.NewReader(synapseAbiStrOther))

	if err != nil {
		panic(err)
	}
}

const (
	Transfer = iota + 1
	AttestMeta
	TransferWithPayload
)

type ParsedVaa struct {
	Version            uint8    `json:"-"`
	GuardianSetIndex   uint32   `json:"-"`
	GuardianSignatures []string `json:"-"`
	Timestamp          uint32
	Nonce              uint32
	EmitterChain       uint16
	EmitterAddress     string
	Sequence           uint64
	ConsistencyLevel   uint8 `json:"-"`
	Payload            string
	Hash               string
}

type TokenTransfer struct {
	PayloadType          uint8
	Amount               *big.Int
	TokenAddress         string
	TokenChain           uint16
	To                   string
	ToChain              uint16
	Fee                  *big.Int
	FromAddress          string
	TokenTransferPayload string
}

// https://github.com/wormhole-foundation/wormhole/blob/048b8834c9/sdk/js/src/vaa/wormhole.ts
func ParseVAA(vm []byte) *ParsedVaa {
	if len(vm) < 6 {
		return nil
	}
	ret := &ParsedVaa{
		Version:            vm[0],
		GuardianSetIndex:   binary.BigEndian.Uint32(vm[1:5]),
		GuardianSignatures: make([]string, 0),
	}
	sigStart := 6
	numSigners := uint(vm[5])
	sigLength := 66
	for i := 0; i < int(numSigners); i++ {
		start := sigStart + i*sigLength
		if len(vm) < start+66 {
			return nil
		}
		ret.GuardianSignatures = append(ret.GuardianSignatures, hexutil.Encode(vm[start+1:start+66]))
	}
	body := vm[sigStart+sigLength*int(numSigners):]
	if len(body) < 51 {
		return nil
	}
	ret.Timestamp = binary.BigEndian.Uint32(body[:4])
	ret.Nonce = binary.BigEndian.Uint32(body[4:8])
	ret.EmitterChain = binary.BigEndian.Uint16(body[8:10])
	ret.EmitterAddress = hexutil.Encode(body[10:42])
	ret.Sequence = binary.BigEndian.Uint64(body[42:50])
	ret.ConsistencyLevel = body[50]
	ret.Payload = hexutil.Encode(body[51:])
	ret.Hash = hexutil.Encode(crypto.Keccak256(body))
	return ret
}

// https://github.com/wormhole-foundation/wormhole/blob/048b8834c9/sdk/js/src/vaa/tokenBridge.ts
func ParseTokenTransferPayload(payload []byte) *TokenTransfer {
	if len(payload) < 133 {
		return nil
	}
	ret := &TokenTransfer{}
	ret.PayloadType = uint8(payload[0])
	if ret.PayloadType != Transfer && ret.PayloadType != TransferWithPayload {
		return nil
	}

	ret.Amount = new(big.Int).SetBytes(payload[1:33])
	ret.TokenAddress = hexutil.Encode(payload[33:65])
	ret.TokenChain = binary.BigEndian.Uint16(payload[65:67])
	ret.To = hexutil.Encode(payload[67:99])
	ret.ToChain = binary.BigEndian.Uint16(payload[99:101])
	if ret.PayloadType == Transfer {
		ret.Fee = new(big.Int).SetBytes(payload[101:133])
	} else {
		ret.FromAddress = hexutil.Encode(payload[101:133])
	}
	ret.TokenTransferPayload = hexutil.Encode(payload[133:])
	return ret
}

func deNormalizeAmount(amount *big.Int, decimals uint8) *big.Int {
	if amount == nil {
		return nil
	}
	if decimals > 8 {
		factor := new(big.Int)
		new(big.Float).SetFloat64(math.Pow10(int(decimals) - 8)).Int(factor)
		amount.Mul(amount, factor)
		return new(big.Int).Set(amount)
	}
	return new(big.Int).Set(amount)
}

func truncateAddress(addr string) string {
	if strings.HasPrefix(addr, "0x000000000000000000000000") {
		addr = "0x" + strings.TrimPrefix(addr, "0x000000000000000000000000")
	}
	return addr
}
