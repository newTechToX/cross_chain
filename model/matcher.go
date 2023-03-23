package model

type Matcher interface {
	Match([]*Data) (Datas, map[string]Datas, error)
	UpdateAnyswapMatchTag(Datas) (int, []*error)
	LastUnmatchId() uint64
	SendMail(string, []*Data)
	Project() string
}
