package aml

type Config struct {
	AmlApiKey string `yaml:"AmlApiKey"`
}

type AddressInfo struct {
	Chain   string
	Address string
	Name    string
	Labels  []string
	Risk    int
}

type RawMsg struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []*ReturnData `json:"data"`
}

type ReturnData struct {
	Address               string                   `json:"address"`
	IsValid               bool                     `json:"is_address_valid"`
	Chain                 string                   `json:"chain"`
	IsContract            bool                     `json:"is_contract"`
	Labels                *Labels                  `json:"labels"`
	CompatibleChainLabels []*CompatibleChainLabels `json:"compatible_chain_labels"`
	Risk                  int                      `json:"risk"`
}
