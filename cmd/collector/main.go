package main

import (
	"app/aggregator"
	"app/config"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/ethereum/go-ethereum/log"
)

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
	go func() {
		fmt.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", *pprofPort), nil))
	}()
	fmt.Println("log level:", lvl.String(), "\nbatch size:", *batchSize, "\npprof port:", *pprofPort, "\nchainbase rate:", *chainbaseRate)
	// log.Root().SetHandler(log.LvlFilterHandler(
	// 	lvl, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	// ))
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler("./error.log", log.TerminalFormat(false))),
		log.LvlFilterHandler(lvl, log.StreamHandler(os.Stderr, log.TerminalFormat(true))),
	))
	ctx, cancel := context.WithCancel(context.Background())
	chainbase.SetupLimit(*chainbaseRate)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		cancel()
	}()
	var cfg config.Config

	config.LoadCfg(&cfg, "../config.yaml")
	srvCtx := svc.NewServiceContext(ctx, &cfg)
	for name := range srvCtx.Config.ChainProviders {
		agg := aggregator.NewAggregator(srvCtx, name)
		go agg.Start()
	}
	<-ctx.Done()
	srvCtx.Wg.Wait()
	fmt.Println("exit")
}
