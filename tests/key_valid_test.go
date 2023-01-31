package tests

import (
	"app/provider"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

func testrovider(chain string, p *provider.Provider) {
	log.Info("validate start", "chain", chain, "len", len(srvCtx.Config.ChainProviders[chain].ApiKeys))
	for i := 0; i < len(srvCtx.Config.ChainProviders[chain].ApiKeys); i++ {
		_, err := p.LatestNumber()
		if err != nil {
			log.Error("failed", "chain", chain, "err", err)
		}
		time.Sleep(time.Second)
	}
}

func TestKeyValid(t *testing.T) {
	var wg sync.WaitGroup
	for chain, p := range srvCtx.Providers.GetAll() {
		wg.Add(1)
		go func(chainKey string, p *provider.Provider) {
			defer wg.Done()
			testrovider(chainKey, p)
		}(chain, p)
	}
	wg.Wait()
	fmt.Println("done")
}
