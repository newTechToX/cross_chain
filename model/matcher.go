package model

type Matcher interface {
	Match([]*Data) (Datas, error)
	UpdateAnyswapMatchTag(Datas) (int, []*error)
	LastUnmatchId() uint64
	Project() string
}
