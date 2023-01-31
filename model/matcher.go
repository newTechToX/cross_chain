package model

type Matcher interface {
	Match([]*Result) (Results, error)
}
