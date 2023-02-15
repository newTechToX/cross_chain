package tests

import (
	"fmt"
	"testing"
	"time"
)

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
