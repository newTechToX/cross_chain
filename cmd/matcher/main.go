package main

import (
	"app/config"
	"app/matcher"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/log"
)

var chainbaseRate = flag.Int("chainbase_limit", 20, "setup chainbase query rate (in Second)")
var Anyswap = flag.Uint64("m", uint64(7121480), "input anyswap startId")
var Across = flag.Uint64("a", uint64(1547700), "input across startId")
var Synapse = flag.Uint64("s", uint64(2), "input synapse startId")

//var Synapse = flag.Uint64("s", uint64(888), "input synapse startId")
//flag.Uint64Var(&r, "synapse", 1234, "help message for flagname")

func main() {
	flag.Parse()
	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}
	var startIds = map[string]uint64{
		"anyswap": *Anyswap,
		"across":  *Across,
		"synapse": *Synapse,
	}

	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))),
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
	m := matcher.NewMatcher(srvCtx, startIds)
	m.Start()
	<-ctx.Done()
	srvCtx.Wg.Wait()
	fmt.Println("exit")
}
