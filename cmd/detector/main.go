package main

import (
	"app/aggregator"
	"app/config"
	"app/detector"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var Anyswap = flag.Uint64("m", uint64(7859426), "input anyswap startId")
var Across = flag.Uint64("a", uint64(1645661), "input across startId")
var logLvl = flag.String("log_level", "info", "set log level")
var batchSize = flag.Int("batch_size", 1000, "set fetch batch size")
var pprofPort = flag.String("pprof", "6060", "set pprof port")
var chainbaseRate = flag.Int("chainbase_limit", 20, "setup chainbase query rate (in Second)")

func main() {
	flag.Parse()
	lvl, err := log.LvlFromString(*logLvl)
	if err != nil {
		panic(err)
	}
	aggregator.BatchSize = uint64(*batchSize)

	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}
	var startIds = map[string]uint64{
		"anyswap": *Anyswap,
		"across":  *Across,
	}
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

	config.LoadCfg(&cfg, "./config.yaml")
	srvCtx := svc.NewServiceContext(ctx, &cfg)
	l := detector.NewDetector(srvCtx, "./logic/txt_config.yaml", startIds)
	l.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	//println("totalFetched ", c)
	<-ctx.Done()
	srvCtx.Wg.Wait()
	fmt.Println("exit")
}
