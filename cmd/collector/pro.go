package main

import (
	"app/aggregator"
	"app/config"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

var Name = flag.String("name", "", "input project name")
var From = flag.Uint64("from", uint64(1547700), "input from blocknumber")
var To = flag.Uint64("to", uint64(1547700), "input to blocknumber")
var Chain = flag.String("chain", "", "input chain")

func main() {
	var logLvl = flag.String("log_level", "info", "set log level")
	var batchSize = flag.Int("batch_size", 1000, "set fetch batch size")
	var pprofPort = flag.String("pprof", "6060", "set pprof port")
	var chainbaseRate = flag.Int("chainbase_limit", 20, "setup chainbase query rate (in Second)")

	flag.Parse()
	lvl, err := log.LvlFromString(*logLvl)
	if err != nil {
		panic(err)
	}
	aggregator.BatchSize = uint64(*batchSize)

	fmt.Println("log level:", lvl.String(), "\nbatch size:", *batchSize, "\npprof port:", *pprofPort, "\nchainbase rate:", *chainbaseRate)
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler("./error.log", log.TerminalFormat(false))),
		log.LvlFilterHandler(lvl, log.StreamHandler(os.Stderr, log.TerminalFormat(true))),
	))
	ctx, _ := context.WithCancel(context.Background())
	chainbase.SetupLimit(*chainbaseRate)
	var cfg config.Config

	config.LoadCfg(&cfg, "./config.yaml")
	srvCtx := svc.NewServiceContext(ctx, &cfg)
	agg := aggregator.NewAggregator(srvCtx, *Chain)
	agg.StartPro(srvCtx, *Name, *From, *To)
	srvCtx.Wg.Wait()
}
