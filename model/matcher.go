package model

type Matcher interface {
	Match(string, []*Data) (Datas, error)
}
