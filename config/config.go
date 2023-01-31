package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Database struct {
	CrosschainDataSource string `yaml:"CrosschainDataSource"`
	AddressLabelSource   string `yaml:"AddressLabelSource"`
	AnyswapSource        string `yaml:"AnyswapSource"`
}

type ChainProvider struct {
	Node            string   `yaml:"Node"`
	ScanUrl         string   `yaml:"ScanUrl"`
	ApiKeys         []string `yaml:"ApiKeys"`
	ChainbaseTable  string   `yaml:"ChainbaseTable"`
	EnableTraceCall bool     `yaml:"EnableTraceCall"`
}

type Config struct {
	Database           Database                  `yaml:"Database"`
	Proxy              string                    `yaml:"Proxy"`
	ChainbaseApiKey    string                    `yaml:"ChainbaseApiKey"`
	EtherscanRateLimit int                       `yaml:"EtherscanRateLimit"`
	ChainProviders     map[string]*ChainProvider `yaml:"ChainProviders"`
	Projects           []string                  `yaml:"Projects"`
}

func LoadCfg[T any](v T, path string) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, v)
	if err != nil {
		panic(err)
	}
}
