package model

type Detector interface {
	DetectOutTx(Datas)
	//DetectOutTx(string, Datas) int
	LastDetectId() uint64
}
