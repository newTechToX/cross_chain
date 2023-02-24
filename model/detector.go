package model

type Detector interface {
	DetectFake(string, Datas) int
	DetectOutTx(string, Datas) int
	LastDetectId() uint64
}
