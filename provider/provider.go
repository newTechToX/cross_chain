package provider

import (
	"app/config"
	"app/model"
	"app/provider/chainbase"
	"app/provider/etherscan"
	"app/provider/geth"
	"fmt"
	"math/big"
)

type Provider struct {
	geth      *geth.GethProvider
	scan      *etherscan.EtherscanProvider
	chainbase *chainbase.Provider
}

func (p *Provider) Call(from, to, input string, value *big.Int, number *big.Int) ([]byte, error) {
	return p.geth.Call(from, to, input, value, number)
}

func (p *Provider) ContinueCall(from, to, input string, value *big.Int, number *big.Int) ([]byte, error) {
	return p.geth.ContinueCall(from, to, input, value, number)
}

func (p *Provider) GetContractFirstInvocation(address string) (uint64, error) {
	val, err := p.scan.GetContractFirstInvocation(address)
	if err == nil {
		return val, nil
	}
	if p.chainbase == nil {
		return 0, err
	}
	return p.chainbase.GetContractFirstCreatedNumber(address)
}

func (p *Provider) LatestNumber() (uint64, error) {
	// val, err := p.geth.LatestNumber()
	// if err == nil {
	// 	return val, nil
	// }
	val, err := p.scan.LatestNumber()
	if err != nil {
		return 0, err
	}
	if val == 0 {
		return 0, fmt.Errorf("invalid latest block")
	}
	return val - 128, nil
}

func (p *Provider) GetLogs(topics0 []string, from, to uint64) (model.Events, error) {
	// return p.geth.GetLogs(addresses, topics0, from, to)
	// must load from chainbase
	return p.scan.GetLogs(topics0, from, to)
	// if p.chainbase == nil {
	// 	return nil, nil
	// }
	// return p.chainbase.GetLogs(addresses, topics0, from, to)
}

func (p *Provider) GetCalls(addresses []string, selectors []string, from, to uint64) ([]*model.Call, error) {
	// return p.scan.GetCalls(addresses, selectors, from, to)
	if p.chainbase == nil {
		return nil, nil
	}
	return p.chainbase.GetCalls(addresses, selectors, from, to)
}

type Providers struct {
	providers map[string]*Provider
}

func NewProviders(cfg *config.Config) *Providers {
	providers := make(map[string]*Provider)
	for chainName, providerCfg := range cfg.ChainProviders {
		gethP := geth.NewGethProvider(chainName, providerCfg.Node)
		scanP := etherscan.NewEtherScanProvider(providerCfg.ScanUrl, providerCfg.ApiKeys, cfg.Proxy, cfg.EtherscanRateLimit)
		providers[chainName] = &Provider{
			geth: gethP,
			scan: scanP,
		}
		if providerCfg.ChainbaseTable != "" {
			providers[chainName].chainbase = chainbase.NewProvider(providerCfg.ChainbaseTable, cfg.ChainbaseApiKey, providerCfg.EnableTraceCall, cfg.Proxy)
		}
	}
	return &Providers{providers: providers}
}

func (p *Providers) Get(chain string) *Provider {
	if chain == "ethereum" {
		chain = "eth"
	}
	if val, ok := p.providers[chain]; ok {
		return val
	}
	return nil
}

func (p *Providers) GetAll() map[string]*Provider {
	return p.providers
}

//  ----------------------

func (p *Provider) GetContractInfo(token string) (*model.ContractInfo, error) {
	info, err := p.chainbase.GetContractInfo(token)
	if info == nil {
		info, err = p.scan.GetContractInfo(token)
	}
	return info, err
}

func (p *Provider) GetLogsWithHash(chain, hash string, topics0 []string) (model.Events, error) {
	events, err := p.chainbase.GetLogsWithHash(chain, hash, topics0)
	return events, err
}

func (p *Provider) GetTraces(chain, hash string) ([]*chainbase.SyChainbaseInfo, error) {
	var res = []*chainbase.SyChainbaseInfo{}
	var err error
	res, err = p.chainbase.GetTraces(chain, hash)
	return res, err
}

func (p *Provider) GetSender(chain, hash string) string {
	sender := p.chainbase.GetSender(chain, hash)
	return sender
}
