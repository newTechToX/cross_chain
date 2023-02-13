package matcher

import (
	"app/dao"
	"app/model"
	"app/utils"
	"fmt"
	"testing"
)

func TestSimpleMatcher(t *testing.T) {
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
}

func TestMatcher_BeginMatch(t *testing.T) {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	tt := 0
	for i := range s {
		go sum(s[i:], c)
		x := <-c
		tt += x
		println(x)
	}
	println(tt)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // 把 sum 发送到通道 c
}
