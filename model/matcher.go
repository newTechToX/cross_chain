package model

type Matcher interface {
	Match(string, []*Data) (Datas, error)
	UpdateAnyswapMatchTag(string, Datas) (int, []*error)
	LastId() uint64
}
