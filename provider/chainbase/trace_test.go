package chainbase

import (
	"app/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/time/rate"
)

func init() {
	SetupLimit(10)
}

func TestExec1(t *testing.T) {
	ret, err := Exec[*Trace]("select * from ethereum.trace_calls where block_number >= 16000000 and block_number < 16000003", "2FtLTBTxc9h7CX3YwBeEkrMlnhc", "")
	fmt.Println(err)
	utils.PrintPretty(ret)
}

func TestExec2(t *testing.T) {
	hash := "0x20ed37da82fd0c2ae9ad0fdd699bce1f40a2ce06e37b06a32001f4ff08cb3433"
	type Sender struct {
		S string `json:"from_address"`
	}
	stmt := fmt.Sprintf("select from_address from %s.transactions where transaction_hash='%s'", "arbitrum", hash)
	ret, err := Exec[*Sender](stmt, "2FtLTBTxc9h7CX3YwBeEkrMlnhc", "")
	if err != nil {
		fmt.Println(err)
	}
	println(ret[0].S)

}

/*func TestExec3(t *testing.T) {
	d := dao.NewDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "select chain, hash from across where safe = 'F' and chain != 'boba'"
	var res []*dao.Data
	err := d.DB().Select(&res, stmt)
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range res {
		stmt = fmt.Sprintf("select block_timestamp from %s.transactions where transaction_hash = '%s'", r.Chain, r.Hash)
		ret, err := Exec[*Trace](stmt, "2FtLTBTxc9h7CX3YwBeEkrMlnhc", "")
		if err != nil {
			fmt.Println(err)
		}
		if len(ret) != 0 {
			st := fmt.Sprintf("update across set ts = '%s' where hash = '%s'", ret[0].Ts, r.Hash)
			_, err := d.DB().Exec(st)
			if err != nil {
				println(err)
			}
		}
	}
	println(len(res))
}*/

func TestGetCalls(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	p := NewProvider("ethereum", "2FtLTBTxc9h7CX3YwBeEkrMlnhc", false, "")
	ret, err := p.GetCalls([]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"}, []string{"0x2e1a7d4d", "0xa9059cbb"}, 16068492, 16068492)
	fmt.Println(err)
	utils.PrintPretty(len(ret))
}

func TestGetLogs(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	p := NewProvider("ethereum", "2FtLTBTxc9h7CX3YwBeEkrMlnhc", true, "http://192.168.3.59:10809")
	// ret, err := p.GetCalls([]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"}, []string{"0x2e1a7d4d", "0xa9059cbb"}, 16068492, 16068492)
	ret, err := p.GetLogs([]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"}, []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65"}, 16068492, 16068493)
	fmt.Println(err)
	utils.PrintPretty(ret)

	ret, err = p.GetLogs([]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"}, nil, 16068492, 16068493)
	fmt.Println(err)
	utils.PrintPretty(ret)
}

func TestLatest(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	p := NewProvider("optimism", "2FtLTBTxc9h7CX3YwBeEkrMlnhc", true, "")
	fmt.Println(p.GetLatestNumber())
}

func TestFirstCall(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	p := NewProvider("ethereum", "2FtLTBTxc9h7CX3YwBeEkrMlnhc", true, "")
	fmt.Println(p.GetContractFirstCreatedNumber("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"))
}

type S struct {
	Id int
}
type SS []*S

func (s SS) Len() int           { return len(s) }
func (s SS) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SS) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func TestSli(t *testing.T) {
	a := []*S{{2}, {1}, {3}}
	sort.Stable(SS(a))
	utils.PrintPretty(a)
	var s []uint64
	ret, err := json.Marshal(s)
	fmt.Println(string(ret), err)

	var qqq []byte
	fmt.Println(string(qqq))
}

func TestRate(t *testing.T) {
	limiter := rate.NewLimiter(5, 1)
	for {
		limiter.Wait(context.Background())
		fmt.Println("111")
	}
}
