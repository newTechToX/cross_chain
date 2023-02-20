package model

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
	"time"
)

/*type Result struct {
	Id          uint64        `db:"id"`
	MatchId     sql.NullInt64 `db:"match_id"`
	Chain       string        `db:"chain"`
	Number      uint64        `db:"number"`
	Ts          time.Time     `db:"ts"`
	Index       uint64        `db:"index"`
	Hash        string        `db:"hash"`
	ActionId    uint64        `db:"action_id"`
	Project     string        `db:"project"`
	Contract    string        `db:"contract"`
	Direction   string        `db:"direction"`
	FromChainId *BigInt       `db:"from_chain_id"`
	FromAddress string        `db:"from_address"`
	ToChainId   *BigInt       `db:"to_chain_id"`
	ToAddress   string        `db:"to_address"`
	Token       string        `db:"token"`
	Amount      *BigInt       `db:"amount"`
	MatchTag    string        `db:"match_tag"`
	Detail      []byte        `db:"detail"`
}

type Results []*Result*/

type MatchedId struct {
	SrcID uint64 `db:"src_id"`
	DstID uint64 `db:"dest_id"`
}

type MatchedIds []*MatchedId

type TokenChain struct {
	Token string `db:"token"`
	Chain string `db:"chain"`
	Block uint64 `db:"block_number"`
}

type TokenChains []*TokenChain

type ContractInfo struct {
	Id       uint64    `db:"id"`
	Project  string    `db:"project"`
	Chain    string    `db:"chain"`
	Address  string    `db:"address"`
	Type     string    `db:"type"`
	Safe     string    `db:"safe"`
	Hash     string    `db:"hash"`
	Ts       time.Time `db:"ts"`
	Number   uint64    `db:"block_number"`
	Deployer string    `db:"deployer"`
}

type ContractInfos []*ContractInfo

type LabelInfo struct {
	NameTag      sql.NullString `db:"name_tag"`
	Labels       sql.NullString `db:"labels"`
	ContractName sql.NullString `db:"contract_name"`
	TokenName    sql.NullString `db:"token_name"`
}

type Data struct {
	Id               uint64         `db:"id"`
	MatchId          sql.NullInt64  `db:"match_id"`
	Chain            string         `db:"chain"`
	Number           uint64         `db:"block_number"`
	Index            uint64         `db:"tx_index"`
	Hash             string         `db:"hash"`
	ActionId         uint64         `db:"log_index"`
	Contract         string         `db:"contract"`
	Direction        string         `db:"direction"`
	FromChainId      *BigInt        `db:"from_chain"`
	FromAddress      string         `db:"from_address"`
	ToChainId        *BigInt        `db:"to_chain"`
	ToAddress        string         `db:"to_address"`
	Token            string         `db:"token"`
	Amount           *BigInt        `db:"amount"`
	MatchTag         string         `db:"match_tag"`
	Profit           []byte         `db:"profit"`
	FromAddressError sql.NullInt64  `db:"from_address_error"`
	ToAddressProfit  sql.NullInt64  `db:"to_address_profit"`
	TokenProfitError sql.NullInt64  `db:"token_profit_error"`
	IsFakeToken      sql.NullInt64  `db:"isfaketoken"`
	Tag              sql.NullString `db:"tag"`
	Project          string         `db:"project"`
}

type Datas []*Data

var CommonFiledNames = []string{
	"id", "chain", "block_number", "tx_index", "hash", "log_index", "contract",
	"direction", "from_chain", "from_address", "to_chain", "to_address", "token",
	"amount", "match_tag", "profit", "from_address_error", "to_address_profit",
	"token_profit_error", "isfaketoken", "tag",
}

var (
	ResultFieldNames = builder.RawFieldNames(&Data{}, true)
	ResultRows       = strings.Join(stringx.Remove(ResultFieldNames, "project"), ",")
	ResultTags       = builder.PostgreSqlJoin(ResultFieldNames)
)
