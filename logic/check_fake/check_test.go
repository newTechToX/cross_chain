package check_fake

import (
	"app/config"
	"app/model"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"os"
	"testing"
)

var cfg config.Config
var srvCtx *svc.ServiceContext

func init() {
	config.LoadCfg(&cfg, "../../config.yaml")
	srvCtx = svc.NewServiceContext(context.Background(), &cfg)
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	chainbase.SetupLimit(10)
}

func TestChecker_IsFakeToken(t *testing.T) {
	id := 7724704
	a := NewFakeChecker(srvCtx, "bsc", "../txt_config.yaml")
	var data model.Datas
	stmt := fmt.Sprintf("select %s from anyswap where id = %d", model.ResultRows, id)
	a.svc.Dao.DB().Select(&data, stmt)
	a.IsFakeToken("anyswap", data[0])
}
