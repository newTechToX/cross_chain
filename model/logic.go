package model

type Detector interface {
	Detect(string, []*Data) (Datas, error)
}
