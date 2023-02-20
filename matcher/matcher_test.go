package matcher

import (
	"app/config"
	"app/dao"
	"app/model"
	"app/provider/chainbase"
	"app/svc"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"os"
	"testing"
)

/*func TestSimpleMatcher(t *testing.T) {
	_dao := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := NewSimpleInMatcher(_dao)
	var results model.Datas
	err := _dao.DB().Select(&results, "select * from common_cross_chain where match_tag = '0x0000000000000000000000000000000000000000000000000000000000000005' and direction = 'in'")
	if err != nil {
		fmt.Println(err)
		return
	}
	shouldUpdates, err := m.Match("", results)
	fmt.Println(err)
	utils.PrintPretty(shouldUpdates)

	err = _dao.Update(shouldUpdates)
	fmt.Println(err)
}*/

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

func TestSimpleInMatcher_UpdateAnyswapMatchTag(t *testing.T) {
	d := dao.NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	id := uint64(352)
	stmt := fmt.Sprintf("select %s from across where id = %d", model.ResultRows, id)
	var data model.Datas
	_ = d.DB().Select(&data, stmt)
	var m = &Matcher{}
	m.svc = srvCtx
	a := NewSimpleInMatcher(srvCtx.ProjectsDao, uint64(6972699))
	m.BeginMatch(id-1, id+1, "across", a)
	println(data[0].MatchTag)
}

func TestNewSimpleInMatcher(t *testing.T) {
	ss := "0x8dceda860bf5d56dce08dbe172dafe6faf79b9ea197b5357fb3939add6b49cb8"
	tt := common.BytesToHash([]byte(ss))
	fmt.Println(tt)
}
