package chainbase

type Number struct {
	Number string `json:"number"`
}

type Log struct {
	Number    string `json:"block_number"`
	Ts        string `json:"block_timestamp"`
	Index     uint64 `json:"transaction_index"`
	Hash      string `json:"transaction_hash"`
	LogIndex  uint64 `json:"log_index"`
	Address   string `json:"address"`
	TopicsCnt uint64 `json:"topics_count"`
	T0        string `json:"topic0"`
	T1        string `json:"topic1"`
	T2        string `json:"topic2"`
	T3        string `json:"topic3"`
	Data      string `json:"data"`
}

type Trace struct {
	Number string `json:"block_number"`
	Ts     string `json:"block_timestamp"`
	Index  uint64 `json:"transaction_index"`
	Hash   string `json:"transaction_hash"`
	From   string `json:"from_address"`
	To     string `json:"to_address"`
	Value  string `json:"value"`
	Input  string `json:"input"`

	// for trace_call
	TraceAddress []uint64 `json:"trace_address"`
}

type Traces []*Trace

func (t Traces) Len() int {
	return len(t)
}

func (t Traces) Less(a, b int) bool {
	if t[a].Number != t[b].Number {
		return t[a].Number < t[b].Number
	}
	if t[a].Index != t[b].Index {
		return t[a].Index < t[b].Index
	}
	return false
}

func (t Traces) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

type Result[T any] struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskId   string `json:"task_id"`
		Result   []T    `json:"result"`
		ErrMsg   string `json:"err_msg"`
		NextPage uint   `json:"next_page"`
	} `json:"data"`
}
