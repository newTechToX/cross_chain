package model

type Detector interface {
	DetectFake(string, Datas) int
	DetectOutTx(string, Datas) int
	LastId() uint64
}
