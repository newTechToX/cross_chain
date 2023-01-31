package etherscan

type EtherscanResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

type GethResponse[T any] struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  T      `json:"result"`
}

type EtherscanEvent struct {
	BlockNumber string   `json:"blockNumber"`
	Timestamp   string   `json:"timeStamp"`
	Index       string   `json:"transactionIndex"`
	LogIndex    string   `json:"logIndex"`
	Address     string   `json:"address"`
	Hash        string   `json:"transactionHash"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
}

type NormalTx struct {
	BlockNumber     string `json:"blockNumber"`
	Timestamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	Index           string `json:"transactionIndex"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Input           string `json:"input"`
	Error           string `json:"isError"`
	ContractAddress string `json:"contractAddress"`
}

type InternalTx struct {
	BlockNumber     string `json:"blockNumber"`
	Timestamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Input           string `json:"input"`
	ContractAddress string `json:"contractAddress"`
	Error           string `json:"isError"`
	Type            string `json:"type"`
}
