package model

type Detector interface {
	DetectFake(string, []*Data) int
	//DetectOutTx(string, []*Data)
}
