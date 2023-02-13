package replay

import (
	"time"
)

type Config struct {
	ReplayAccessToken string `yaml:"ReplayAccessToken"`
}

type RawMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type T struct {
	Txns []*SimulatedTxn `json:"txns"`
}

type SimulatedTxn struct {
	ChainID        int                    `json:"chainID"`
	Block          uint64                 `json:"block"`
	Hash           string                 `json:"hash"`
	From           string                 `json:"from"`
	To             string                 `json:"to"`
	Contract       string                 `json:"contract,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
	Input          string                 `json:"input"` // 拼接了 selector 的 call parameter
	GasLimit       int                    `json:"gasLimit"`
	GasPrice       string                 `json:"gasPrice"` // Wei
	GasUsed        int                    `json:"gasUsed"`
	BaseFee        string                 `json:"baseFee,omitempty"`     // Wei
	MaxFee         string                 `json:"maxFee,omitempty"`      // 不一定有(有则返回, 没有则不返回) Wei
	PriorityFee    string                 `json:"priorityFee,omitempty"` // 不一定有 Wei
	TransactionFee string                 `json:"transactionFee"`        // Wei
	Nonce          int                    `json:"nonce"`
	Position       int                    `json:"position"`
	Type           string                 `json:"type"`  // Legacy | AccessList | EIP-1559
	Value          string                 `json:"value"` // Wei
	Status         bool                   `json:"status"`
	Error          string                 `json:"error"`
	InternalTxns   []*InternalTransaction `json:"internalTxns"`
	Events         []*Event               `json:"events"`
	BalanceChanges []*SimAccountBalance   `json:"balanceChanges"`
}

type InternalTransaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Contract string `json:"contract,omitempty"`
	Input    string `json:"input"` // 拼接了 selector 的 call parameter
	Output   string `json:"output"`
	GasLimit int    `json:"gasLimit"`
	GasUsed  int    `json:"gasUsed"`
	Status   bool   `json:"status"`
	Error    string `json:"error"`
	Type     string `json:"type"` // CALL | DELEGATECALL
	Value    string `json:"value"`
}

type Event struct {
	Address string   `json:"address"`
	Data    string   `json:"data"`
	Topics  []string `json:"topics"`
}

type SimAccountBalance struct {
	Account string      `json:"account"`
	Assets  []*SimAsset `json:"assets"`
}

type SimAsset struct {
	Address  string `json:"address"`
	Amount   string `json:"amount"`          // raw amount
	Decimals int    `json:"decimals"`        // -1 标识缺少 decimals 信息
	Value    string `json:"value,omitempty"` // 空字符串标识缺少价格信息
}
