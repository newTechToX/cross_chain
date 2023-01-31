package model

import (
	"math/big"
	"time"
)

type EventConfig struct {
	Addresses []string
	Topics0   []string
}

type MsgConfig struct {
	Addresses []string
	Selector  []string
}

type Event struct {
	Number  uint64    `json:"number"`
	Ts      time.Time `json:"ts"`
	Index   uint64    `json:"index"`
	Hash    string    `json:"hash"`
	Id      uint64    `json:"id"`
	Address string    `json:"address"`
	Topics  []string  `json:"topics"`
	Data    string    `json:"data"`
}

type Events []*Event

func (e Events) Len() int {
	return len(e)
}

func (e Events) Less(i, j int) bool {
	if e[i].Number != e[j].Number {
		return e[i].Number < e[j].Number
	} else if e[i].Index != e[j].Index {
		return e[i].Index < e[j].Index
	}
	return e[i].Id < e[j].Id
}

func (e Events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type Call struct {
	Number uint64    `json:"number"`
	Ts     time.Time `json:"ts"`
	Index  uint64    `json:"index"`
	Hash   string    `json:"hash"`
	Id     uint64    `json:"id"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Input  string    `json:"input"`
	Value  *big.Int  `json:"value"`
}
