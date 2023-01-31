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

func main() {
	flag.Parse()
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
	m := matcher.NewMatcher(srvCtx)
	m.Start()
	<-ctx.Done()
	srvCtx.Wg.Wait()
	fmt.Println("exit")
}
