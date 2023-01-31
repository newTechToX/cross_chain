package matcher

import (
	"app/dao"
	"app/model"
	"app/utils"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/log"
)

type SimpleInMatcher struct {
	dao *dao.Dao
}

var _ model.Matcher = &SimpleInMatcher{}

func NewSimpleInMatcher(dao *dao.Dao) *SimpleInMatcher {
	return &SimpleInMatcher{dao: dao}
}

// match cross-out txs with cross-in txs, the inputs should be cross-in
// src => dst, match_tag equals, to_chain = from_chain
// inputs: cross-in txs
// require: to_chain_id in cross-out must exist
// matched: match_tags equal
func (m *SimpleInMatcher) Match(crossIns []*model.Result) (shouldUpdates model.Results, err error) {
	for _, crossIn := range crossIns {
		if crossIn.Direction != model.InDirection {
			log.Warn("matching should not input cross-out")
			continue
		}
		var pending model.Results
		stmt := fmt.Sprintf("select * from %s where match_tag = $1 and project = $2 and direction = '%s' and to_chain_id = $3", m.dao.Table(), model.OutDirection)
		err := m.dao.DB().Select(&pending, stmt, crossIn.MatchTag, crossIn.Project, utils.GetChainId(crossIn.Chain).String())
		if err != nil {
			return nil, err
		}
		if len(pending) == 0 {
			continue
		}
		// if len(pending) > 1 {
		// 	log.Error("multi matched", "src", crossIn.Hash)
		// }
		valid := make(model.Results, 0)
		for _, counterparty := range pending {
			if !isMatched(counterparty, crossIn) {
				continue
			}
			valid = append(valid, counterparty)
			fillEmptyFields(counterparty, crossIn)
		}
		if len(valid) > 1 {
			log.Error("multi matched", "src", crossIn.Hash, "chain", crossIn.Chain, "project", crossIn.Project)
		} else {
			shouldUpdates = append(shouldUpdates, crossIn)
			shouldUpdates = append(shouldUpdates, valid...)
		}
	}
	return
}

func isMatched(out, in *model.Result) bool {
	if out.ToChainId != nil {
		if (*big.Int)(out.ToChainId).Cmp(utils.GetChainId(in.Chain)) != 0 {
			return false
		}
	}
	if in.FromChainId != nil {
		if (*big.Int)(in.FromChainId).Cmp(utils.GetChainId(out.Chain)) != 0 {
			return false
		}
	}
	if out.FromAddress != "" && in.FromAddress != "" && out.FromAddress != in.FromAddress {
		return false
	}
	if out.ToAddress != "" && in.ToAddress != "" && out.ToAddress != in.ToAddress {
		return false
	}

	return true
}

func fillEmptyFields(out, in *model.Result) {
	if out == nil || in == nil || out.Direction != model.OutDirection || in.Direction != model.InDirection {
		log.Error("invalid match pair")
		return
	}
	in.MatchId = sql.NullInt64{Int64: int64(out.Id), Valid: true}
	out.MatchId = sql.NullInt64{Int64: int64(in.Id), Valid: true}
	// fill empty in cross-in
	if in.FromChainId == nil {
		in.FromChainId = (*model.BigInt)(new(big.Int).Set(utils.GetChainId(out.Chain)))
	}
	if in.FromAddress == "" {
		in.FromAddress = out.FromAddress
	}
	if in.ToChainId == nil {
		in.ToChainId = (*model.BigInt)(new(big.Int).Set(utils.GetChainId(in.Chain)))
	}
	if in.ToAddress == "" {
		in.ToAddress = out.ToAddress
	}
	//fill empty in cross-out

	if out.FromChainId == nil {
		out.FromChainId = (*model.BigInt)(new(big.Int).Set((utils.GetChainId(out.Chain))))
	}
	if out.FromAddress == "" {
		out.FromAddress = in.FromAddress
	}
	if out.ToChainId == nil {
		out.ToChainId = (*model.BigInt)(new(big.Int).Set(utils.GetChainId(in.Chain)))
	}
	if out.ToAddress == "" {
		out.ToAddress = in.ToAddress
	}
}
