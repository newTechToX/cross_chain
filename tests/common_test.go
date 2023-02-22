package tests

import (
	"app/dao"
	"app/model"
	"fmt"
	"regexp"
	"testing"
	"time"
)

var d = dao.NewAnyDao("postgres://xiaohui_hu:xiaohui_hu_blocksec888@192.168.3.155:8888/cross_chain?sslmode=disable")

func TestSize(t *testing.T) {
	total := 17
	size := total / 3
	i := 0
	for ; i < total-2*size; i = i + size {
		go pp(i, size)
	}
	pp(i, total-i)
}

func pp(i, size int) {
	for j := 0; j < size; j++ {
		println(i)
		i++
		time.Sleep(1 * time.Second)
	}
	println("----------")
}

func TestMatcher_BeginMatch(t *testing.T) {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan []int)
	//d := make(chan []int)
	tt := make(map[int][]int)
	for i := range s {
		go sum(i, s[i:], c)
		x := <-c
		for k, v := range x {
			tt[k] = append(tt[k], v)
		}
	}
	go sum(0, s, c)
	x := <-c
	for k, v := range x {
		tt[k] = append(tt[k], v)
	}
	for k, v := range tt {
		fmt.Println(k, v)
	}
}

func sum(i int, s []int, c chan []int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- []int{sum, i} // 把 sum 发送到通道 c
}

func TestMatcher_Timer(t *testing.T) {
	tim := time.NewTimer(2 * time.Second)
	LAST := 0
	last := 0
	cnt := 0
	for {
		select {
		case <-tim.C:
			cnt += 1
			last += 2
			println("last: ", last)
		}
		if cnt >= 5 {
			last = LAST
		}
		tim.Reset(5 * time.Second)
	}
}

func Test_If_all_srcTxHash_changed(t *testing.T) {
	type Y struct {
		Chain string `db:"chain"`
		Id    uint64 `db:"m"`
	}
	var y []*Y
	stmt := "select chain, min(block_number) as m from anyswap where id > 5654037 and direction = 'in' and (from_chain = '1' or from_chain = '10' or from_chain = '56' or from_chain ='137' or from_chain ='250' or from_chain ='42161' or from_chain = '43114') and match_id is null group by chain"
	err := d.DB().Select(&y, stmt)
	if err != nil {
		fmt.Println("1", err)
	} else {
		for _, r := range y {
			stmt = fmt.Sprintf("select contract from anyswap where chain = '%s' and block_number > %d group by contract", r.Chain, r.Id)
			var datas []*string
			err = d.DB().Select(&datas, stmt)
			if err != nil {
				fmt.Println("2", err)
			}
		}
	}
}

func TestRehix(t *testing.T) {
	var isStringAlphabetic = regexp.MustCompile(`^[0-9a-z]+$`).MatchString
	s := "0xa5722bb24e31b6b4b710183c6ae4518613645aaf"
	sa := "0xA5722bb24e31b6b4b710183c6ae4518613645aaf"
	if isStringAlphabetic(s) {
		println("okok")
	} else if isStringAlphabetic(sa) {
		println("sa")
	}
}

func TestGetArbitrum(t *testing.T) {
	var isStringAlphabetic = regexp.MustCompile(`^[0-9a-z]+$`).MatchString
	stmt := "select * from anyswap where direction = 'out' and chain = 'bsc' and id > 6972698 order by id desc"
	var datas model.Datas
	err := d.DB().Select(&datas, stmt)
	if err != nil {
		fmt.Println(err)
	}
	println(len(datas))
	cnt := 0
	for _, d := range datas {
		if !isStringAlphabetic(d.ToAddress) {
			println(d.Id)
			cnt++
		}
	}
	println(cnt)
}

func TestDeleteDuplicates(t *testing.T) {
	chains := []string{"bsc", "fantom", "arbitrum", "optimism", "ethereum", "avalanche"}
	for _, chain := range chains {
		go deleteDuplicates(chain)
	}
	deleteDuplicates("polygon")
}

// 删除重复拉取的数据
func deleteDuplicates(chain string) {
	Web3QueryStartBlock := map[string]uint64{
		"ethereum":  12000000,
		"bsc":       25545001,
		"polygon":   15000000,
		"fantom":    2000000,
		"arbitrum":  900,
		"avalanche": 2400000,
		"optimism":  3400000,
	}
	Web3QueryEndBlock := map[string]uint64{
		"ethereum":  16661034,
		"polygon":   39460673,
		"bsc":       25795917,
		"avalanche": 26449812,
		"fantom":    56086014,
		"arbitrum":  62348616,
		"optimism":  75281535,
	}
	stmt := fmt.Sprintf("select * from anyswap where direction = 'out' and chain = '%s' and id > 7400000 and block_number > %d and block_number < %d",
		chain, Web3QueryStartBlock[chain], Web3QueryEndBlock[chain])
	var datas = model.Datas{}
	er := d.DB().Select(&datas, stmt)
	if er != nil {
		fmt.Println(er)
		return
	}
	println(chain, len(datas))

	cnt := 0
	for _, data := range datas {
		stmt = fmt.Sprintf("select * from anyswap where hash = '%s' and log_index = %d and tx_index = %d and id != %d", data.Hash, data.ActionId, data.Index, data.Id)
		var dup model.Datas
		er = d.DB().Select(&dup, stmt)
		if er != nil {
			fmt.Println(er)
			continue
		}
		for _, dd := range dup {
			cnt++
			stmt = fmt.Sprintf("delete from anyswap where id = %d", dd.Id)
			_, err := d.DB().Exec(stmt)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	println(chain, "deleted ", cnt)
}
