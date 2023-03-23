package tests

import (
	"app/dao"
	"app/model"
	"app/utils"
	"database/sql"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"regexp"
	"sync"
	"sync/atomic"
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
		stmt = fmt.Sprintf("select * from anyswap where hash = '%s' and log_index = %d and tx_index = %d and id != %d", data.Hash, data.LogIndex, data.TxIndex, data.Id)
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

func Test_Copy_slice(t *testing.T) {
	a := []int{1, 2, 3}
	for i := range a[1:] {
		a = a[:i+1+copy(a[i+1:], a[i+2:])]
	}
	fmt.Println(a)
}

func Test_go(t *testing.T) {
	t1 := time.Now()

	list := []int{1, 2, 5}
	res := collect(list)
	t2 := time.Now()
	fmt.Println(t2.Sub(t1).String())
	fmt.Println(res)
}

func fu(list []int, response chan int, limiter chan bool, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {

	defer wg.Done()
	for _, i := range list {
		time.Sleep(2 * time.Second)
		println(i)
	}
	response <- list[0]
	bar.Add(1)
	<-limiter
}

func collect(urls []int) []int {
	var result []int
	//var donCh chan struct{}
	bar := utils.Bar(3, "do")
	wg := &sync.WaitGroup{}
	// 控制并发数为10
	limiter := make(chan bool, 10)
	defer close(limiter)
	responseChannel := make(chan int, 20)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}
	// 启动读取结果的控制器
	go func() {
		// wgResponse计数器+1
		wgResponse.Add(1)
		// 读取结果
		for response := range responseChannel {
			// 处理结果
			result = append(result, response)
		}
		// 当 responseChannel被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
		wgResponse.Done()
	}()
	for i := range urls {
		// 计数器+1
		wg.Add(1)
		limiter <- true
		// 这里在启动goroutine时, 将用来收集结果的局部变量channel也传递进去
		go fu(urls[i:], responseChannel, limiter, wg, bar)
	}

	// 等待所以协程执行完毕
	wg.Wait() // 当计数器为0时, 不再阻塞
	fmt.Println("所有协程已执行完毕")

	// 关闭接收结果channel
	close(responseChannel)

	// 等待wgResponse的计数器归零
	wgResponse.Wait()

	// 返回聚合后结果
	return result
}

func TestSqlNULLint(t *testing.T) {
	d := model.Data{}
	d.IsFakeToken.Scan(1)
	if d.IsFakeToken.Valid && d.IsFakeToken.Int64 == 1 {
		println("pk")
	}
}

type Sli []int

func TestSlice(t *testing.T) {
	var c = Sli{
		1, 2, 3, 4, 5,
	}
	c = deleteSlice(c, []int{1, 2, 6, 5})
	fmt.Println(c)
}

func deleteSlice(a Sli, ids []int) Sli {
	k := -1
	b := []int{}
	for i, id := range ids {
		for j := 0; j < len(a); j++ {
			aa := a[j]
			if id == aa {
				b = append(b, ids[k+1:i]...)
				k = i
				fmt.Println(b)
				break
			}
		}
	}
	b = append(b, ids[k+1:]...)
	return b
}

func TestMap(t *testing.T) {
	a := map[int][]int{
		1: {1, 2}, 4: {2, 3},
	}
	a[5] = append(a[5], 9)
	fmt.Println(a[5])
}

func TestChan(t *testing.T) {
	num := map[int]int{}
	list := []int{1, 2, 3, 4, 5}
	//channels := make(chan int)
	aa := 0
	Wg := &sync.WaitGroup{}
	t1 := time.Now()
	for _, i := range list {
		go testChan(i, num, Wg)
		//aa += <-channels
	}
	testChan(1, num, Wg)
	t2 := time.Now()
	fmt.Println(t2.Sub(t1).String(), aa)
	fmt.Println(num)
}
func testChan(a int, num map[int]int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for i := 0; i < a; i++ {
		time.Sleep(time.Second)
	}
	num[a] = a + 1
}

func TestContinur(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		switch a[i] {
		case 3:
			if a[i] == 4 {
				continue
			}
			println(a[i])
		default:
			println(a[i])
		}
	}
}

func f(i int) int {
	time.Sleep(time.Duration(i) * time.Second) // 模拟函数执行需要一定时间
	return i * 2                               // 返回i的两倍
}

func TestGPTch(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int, 10) // 设置通道的缓冲区大小为10

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result := f(i)
			ch <- result // 将结果发送到通道中
		}(i)
	}

	wg.Wait()

	close(ch) // 关闭通道，表示所有协程已经完成

	sum := 0
	for result := range ch {
		sum += result // 从通道中获取结果，并累加到sum变量中
	}

	fmt.Println("Sum:", sum)
}

type SafeWaitGroup struct {
	wg  sync.WaitGroup
	cnt int32
}

func (swg *SafeWaitGroup) Add(delta int) {
	atomic.AddInt32(&swg.cnt, int32(delta))
	swg.wg.Add(delta)
}

func (swg *SafeWaitGroup) Done() {
	cnt := atomic.AddInt32(&swg.cnt, -1)
	if cnt < 0 {
		panic("negative WaitGroup counter")
	}
	swg.wg.Done()
}

func (swg *SafeWaitGroup) Wait() {
	swg.wg.Wait()
}

func main() {
	var swg SafeWaitGroup

	for i := 0; i < 10; i++ {
		swg.Add(1)
		go func(i int) {
			defer swg.Done()
			fmt.Println(i * 2)
		}(i)
	}

	swg.Wait()
}

func TestScan(t *testing.T) {
	var a sql.NullString
	if !a.Valid || a.String == "" {
		_ = a.Scan("iuouo")
	}
	println(a.String)
}

type A struct {
	a int
}
type B struct {
	b int
	c []*A
}

func TestAppend(t *testing.T) {
	s := make(map[int]*B, 1)
	for i := 0; i < 3; i++ {
		te(s, i)
	}

	for v, k := range s {
		fmt.Println(v, k)
	}
	return
}

func te(s map[int]*B, i int) {
	if _, ok := s[i]; !ok {
		s[i] = &B{
			b: i,
			c: []*A{},
		}
	}
	s[i].c = append(s[i].c, &A{a: i + 1})
}
