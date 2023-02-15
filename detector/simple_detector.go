package detector

import (
	"app/logic"
	"app/logic/check_fake"
	"app/matcher"
	"app/model"
	"app/svc"
	"github.com/ethereum/go-ethereum/log"
	log2 "log"
)

type SimpleOutDetector struct {
	svc     *svc.ServiceContext
	logic   *logic.Logic
	matcher *matcher.Matcher
}

var _ model.Detector = &SimpleOutDetector{}

func NewSimpleOutDetector(svc *svc.ServiceContext, chain string, config_path string) *SimpleOutDetector {
	return &SimpleOutDetector{
		svc:     svc,
		logic:   logic.NewLogic(svc, chain, config_path),
		matcher: matcher.NewMatcher(svc),
	}
}

// OutDetector的 Detect 用于检测所有tx的fake token & chainID，将没有match的做二次检测
// 对于fake chainID的检查还没做完

func (a *SimpleOutDetector) DetectFake(project string, datas model.Datas) int {
	var unsafe_chan = make(chan map[int]model.Datas)
	var size, i = len(datas) / 20, 0
	var res = make(map[int]model.Datas)

	//协程处理datas
	for ; i < len(datas)-2*size; i = i + size {
		go a.logic.CheckFake(project, datas[i:i+size], unsafe_chan)
		unsafe_map := <-unsafe_chan
		for tag, dts := range unsafe_map {
			res[tag] = append(res[tag], dts...)
		}
	}
	go a.logic.CheckFake(project, datas[i:], unsafe_chan)
	unsafe_map := <-unsafe_chan
	for tag, dts := range unsafe_map {
		res[tag] = append(res[tag], dts...)
	}

	if len(res) != 0 {
		for tag, f := range res {
			for _, d := range f {
				log2.SetPrefix("DetectFake()")
				log.Error("fake token detected", "project", project, "tag", tag, "token", d.Token, "hash", d.Hash)
			}

		}
	}
	return len(res[check_fake.FAKE_TOKEN]) + len(res[check_fake.NULL_IN_DATABASE])
}

func (a *SimpleOutDetector) DetectOutTx(project string, datas model.Datas) int {
	var unsafe_chan = make(chan map[int]model.Datas)
	var size, i = len(datas) / 20, 0
	var res = make(map[int]model.Datas)

	//协程处理datas
	for ; i < len(datas)-2*size; i = i + size {
		go a.logic.CheckFake(project, datas[i:i+size], unsafe_chan)
		unsafe_map := <-unsafe_chan
		for tag, dts := range unsafe_map {
			res[tag] = append(res[tag], dts...)
		}
	}
	go a.logic.CheckFake(project, datas[i:], unsafe_chan)
	unsafe_map := <-unsafe_chan
	for tag, dts := range unsafe_map {
		res[tag] = append(res[tag], dts...)
	}
	return 0
}
