package model

type Colletcor interface {
	Name() string
	Contracts(chain string) []string
}

type EventCollector interface {
	Colletcor
	Topics0(chain string) []string
	Extract(chain string, events Events) Datas
}

type MsgCollector interface {
	Colletcor
	Selectors(chain string) []string
	Extract(chain string, msgs []*Call) Datas
}
