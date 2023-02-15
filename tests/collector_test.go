package tests

import (
	"app/config"
	"app/cross_chain/anyswap"
	"app/cross_chain/celer_bridge"
	"app/cross_chain/hop"
	renbridge "app/cross_chain/ren_bridge"
	"app/cross_chain/stargate"
	"app/cross_chain/wormhole"
	"app/model"
	"app/provider/chainbase"
	"app/svc"
	"app/utils"
	"context"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

var cfg config.Config
var srvCtx *svc.ServiceContext

func init() {
	config.LoadCfg(&cfg, "../config.yaml")
	srvCtx = svc.NewServiceContext(context.Background(), &cfg)
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	chainbase.SetupLimit(10)
}

func TestAnyswap(t *testing.T) {
	c := anyswap.NewAnyswapCollector(srvCtx)
	addrs := c.Contracts("fantom")
	if addrs == nil {
		return
	}
	topics0 := c.Topics0("fantom")

	p := srvCtx.Providers.Get("fantom")

	calls, err := p.GetLogs(addrs, topics0, 35860281, 35860281)
	if err != nil {
		return
	}
	// utils.PrintPretty(addrs)
	results := c.Extract("eth", calls)
	utils.PrintPretty(results)
}

func TestAnyswapUnderlying(t *testing.T) {
	any := anyswap.NewAnyswapCollector(srvCtx)
	fmt.Println(any.GetUnderlying("eth", "0x22648C12Acd87912ea1710357b1302c6a4154ebc"))
}

func TestCBridge(t *testing.T) {
	c := celer_bridge.NewCBridgeCollector()
	fmt.Println(c.Name())
	p := srvCtx.Providers.Get("eth")
	es, err := p.GetLogs(
		[]string{"0xb37d31b2a74029b5951a2778f959282e2d518595"},
		[]string{"0x15d2eeefbe4963b5b2178f239ddcc730dda55f1c23c22efb79ded0eb854ac789"},
		16045356, 16045356)
	fmt.Println(err)
	ret := c.Extract("eth", es)
	fmt.Println(ret)
}

func TestStargate(t *testing.T) {
	c := stargate.NewStargateCollector(srvCtx)
	// for bsc
	addrs := c.Contracts("bsc")
	if addrs == nil {
		return
	}
	topics0 := c.Topics0("bsc")
	p := srvCtx.Providers.Get("bsc")
	events, err := p.GetLogs(addrs, topics0, 23499162, 23499162)
	if err != nil {
		return
	}
	sort.Sort(events)
	results := c.Extract("bsc", events)
	utils.PrintPretty(results)

	// for eth
	addrs = c.Contracts("eth")
	if addrs == nil {
		return
	}
	topics0 = c.Topics0("eth")
	p = srvCtx.Providers.Get("eth")
	events, err = p.GetLogs(addrs, topics0, 16081998, 16081998)
	if err != nil {
		return
	}
	sort.Sort(events)
	results = c.Extract("eth", events)
	utils.PrintPretty(results)

	// for bsc
	addrs = c.Contracts("bsc")
	if addrs == nil {
		return
	}
	topics0 = c.Topics0("bsc")
	p = srvCtx.Providers.Get("bsc")
	events, err = p.GetLogs(addrs, topics0, 23508697, 23508697)
	if err != nil {
		return
	}
	sort.Sort(events)
	results = c.Extract("bsc", events)
	utils.PrintPretty(results)
}

func TestHop(t *testing.T) {
	c := hop.NewHopCollector()
	events := make(model.Events, 0)
	addrs := c.Contracts("polygon")
	if addrs == nil {
		return
	}
	topics0 := c.Topics0("polygon")

	p := srvCtx.Providers.Get("polygon")

	ret, err := p.GetLogs(addrs, topics0, 37111294, 37111294)
	if err != nil {
		return
	}
	events = append(events, ret...)
	sort.Sort(events)
	results := c.Extract("polygon", events)
	println(results)
}

func TestRen(t *testing.T) {
	c := renbridge.NewRenbridgeCollector()

	addrs := c.Contracts("eth")
	if addrs == nil {
		return
	}
	topics0 := c.Selectors("eth")

	p := srvCtx.Providers.Get("eth")

	ret, err := p.GetCalls(addrs, topics0, 16041412, 16041414)
	if err != nil {
		return
	}
	results := c.Extract("eth", ret)
	utils.PrintPretty(results)
}

func TestStargateFindPairs(t *testing.T) {
	c := stargate.NewStargateCollector(srvCtx)

	addrs := c.Contracts("eth")
	if addrs == nil {
		return
	}
	// topics0 := c.Topics0("eth")

	p := srvCtx.Providers.Get("eth")

	ret, err := p.GetLogs(nil, nil, 16083387, 16083387)
	if err != nil {
		return
	}
	pairs := stargate.FindParis(ret, stargate.Swap, stargate.SendMsg)
	utils.PrintPretty(pairs)
	// results := c.Extract("eth", ret)
	// utils.PrintPretty(results)
}

func TestStargateToken(t *testing.T) {
	c := stargate.NewStargateCollector(srvCtx)

	addrs := c.Contracts("eth")
	if addrs == nil {
		return
	}
	// topics0 := c.Topics0("eth")

	fmt.Println(c.GetPoolToken("eth", "0x101816545f6bd2b1076434b54383a1e633390a2e"))

	fmt.Println(c.GetPoolConvertRate("bsc", "0x98a5737749490856b401db5dc27f522fc314a4e1"))
	// results := c.Extract("eth", ret)
	// utils.PrintPretty(results)
}

func TestStargateAbi(t *testing.T) {
	ret, err := stargate.DecodePacketReceivedData("0x00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000ae7fc8e30c8877d723ea5ff5919eda6461c8a75d393c15c9f882d2289c41a361e5000000000000000000000000000000000000000000000000000000000000000146694340fc020c5e6b96567843da2df01b2ce1eb6000000000000000000000000")
	val0, ok := ret[0].([]byte)
	if ok {
		fmt.Println(hexutil.Encode(val0))
	}
	var1, ok := ret[1].(uint64)
	if ok {
		fmt.Println(var1)
	}
	fmt.Println(err, len(ret))
}

func TestWormwhole(t *testing.T) {
	c := wormhole.NewWormHoleCollector(srvCtx)
	addrs := c.Contracts("eth")
	if addrs == nil {
		return
	}
	topics0 := c.Selectors("eth")

	p := srvCtx.Providers.Get("eth")

	calls, err := p.GetCalls(addrs, topics0, 16130849, 16132849)
	if err != nil {
		return
	}
	// utils.PrintPretty(addrs)
	results := c.Extract("eth", calls)
	utils.PrintPretty(results)
}
