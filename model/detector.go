package model

type Detector interface {
	DetectOutTx(Datas) int
	//DetectOutTx(string, Datas) int
	LastDetectId() uint64
}
