package matcher

import (
	"app/dao"
	"app/model"
	"app/utils"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

type SimpleInMatcher struct {
	project  string
	dao      *dao.Dao
	start_id uint64
}

var _ model.Matcher = &SimpleInMatcher{}

func NewSimpleInMatcher(project string, dao *dao.Dao, start_id uint64) *SimpleInMatcher {
	return &SimpleInMatcher{
		project:  project,
		dao:      dao,
		start_id: start_id,
	}
}

func (a *SimpleInMatcher) LastUnmatchId() uint64 {
	stmt := fmt.Sprintf("select min(id) from %s where direction = 'in' and id >= %d and match_id is null and from_chain in (1, 10, 56, 137, 250, 42161, 43114)", a.project, a.start_id)
	type ID struct {
		Id uint64 `db:"min"`
	}
	var id = ID{a.start_id}
	if err := a.dao.DB().Get(&id, stmt); err != nil {
		log.Warn("failed to get unmatchId", "project", a.project, "ERROR", err)
	} else {
		a.start_id = id.Id
	}
	return a.start_id
}

func (a *SimpleInMatcher) Project() string {
	return a.project
}

// match cross-out txs with cross-in txs, the inputs should be cross-in
// src => dst, match_tag equals, to_chain = from_chain
// inputs: cross-in txs
// require: to_chain_id in cross-out must exist
// matched: match_tags equal

func (a *SimpleInMatcher) Match(crossIns []*model.Data) (shouldUpdates model.Datas, err error) {
	for _, crossIn := range crossIns {
		if crossIn.Direction != model.InDirection {
			log.Warn("matching should not input cross-out")
			continue
		}
		var pending model.Datas
		var stmt string
		var err error

		switch a.project {
		case "across":
			stmt = fmt.Sprintf("select %s from %s where match_tag = $1 and direction = '%s' and to_chain = $2 and from_address = $3 and to_address = $4 and amount = $5", model.ResultRows, a.project, model.OutDirection)
			err = a.dao.DB().Select(&pending, stmt, crossIn.MatchTag, utils.GetChainId(crossIn.Chain).String(), crossIn.FromAddress, crossIn.ToAddress, crossIn.Amount.String())
		default:
			stmt = fmt.Sprintf("select %s from %s where match_tag = $1 and direction = '%s' and to_chain = $2", model.ResultRows, a.project, model.OutDirection)
			err = a.dao.DB().Select(&pending, stmt, crossIn.MatchTag, utils.GetChainId(crossIn.Chain).String())
		}
		if err != nil {
			return nil, err
		}
		if len(pending) == 0 {
			continue
		}
		// if len(pending) > 1 {
		// 	log.Error("multi matched", "src", crossIn.Hash)
		// }
		valid_ := make(model.Datas, 0)
		multi := 0
		var dups = model.Datas{crossIn}
		for _, counterparty := range pending {
			if !isMatched(counterparty, crossIn) {
				continue
			}

			if counterparty.MatchId.Valid {
				multi = 1
				//说明已经match过，但有可能是数据重复的原因
				stmt = fmt.Sprintf("select %s from %s where id = %d", model.ResultRows, a.project, counterparty.MatchId.Int64)
				var dup model.Data
				if err = a.dao.DB().Get(&dup, stmt); err != nil {
					fmt.Println(err)
				} else if isDuplicated(&dup, crossIn) {
					multi = 2
					err = a.dao.Delete(a.Project(), crossIn.Id)
				} else {
					dups = append(dups, &dup)
					dups = append(dups, counterparty)
				}
			} else {
				//如果没有匹配的，就添加入更新列表
				valid_ = append(valid_, counterparty)
				fillEmptyFields(counterparty, crossIn)
			}
		}
		if len(valid_) == 0 {
			if multi == 0 {
				a.SendMail("UNMATCH", dups)
				log.Error("unmatch", "src", crossIn.Hash, "chain", crossIn.Chain, "project", a.project)
			} else if multi == 1 {
				a.SendMail("MULTI MATCHED IN", dups)
				log.Error("in tx multi matched", "src", crossIn.Hash, "chain", crossIn.Chain, "project", a.project)
			}
			continue
		}

		valid := model.Datas{valid_[0]}
		if len(valid_) > 1 {
			for _, data := range valid_[1:] {
				if isDuplicated(data, valid_[0]) {
					err = a.dao.Delete(a.Project(), data.Id)
				} else {
					valid = append(valid, data)
				}
			}
		}
		if len(valid) > 1 {
			dups = append(dups, valid...)
			a.SendMail("MULTI MATCHED OUT", dups)
			log.Error("out tx multi matched", "src", crossIn.Hash, "chain", crossIn.Chain, "project", a.project)

		}
		if len(valid) >= 1 {
			shouldUpdates = append(shouldUpdates, crossIn)
			shouldUpdates = append(shouldUpdates, valid...)
		}
	}
	return
}

func isMatched(out, in *model.Data) bool {
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
	if out.FromAddress != "" && in.FromAddress != "" && strings.ToLower(out.FromAddress) != strings.ToLower(in.FromAddress) {
		return false
	}
	if out.ToAddress != "" && in.ToAddress != "" && strings.ToLower(out.ToAddress) != strings.ToLower(in.ToAddress) {
		return false
	}

	return true
}

func fillEmptyFields(out, in *model.Data) {
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

func (a *SimpleInMatcher) UpdateAnyswapMatchTag(crossIns model.Datas) (cnt int, errs []*error) {
	shouldUpdates, errs := updateAnyswapMatchTag(crossIns)
	cnt = a.dao.UpdateAnyswapMatchTag(a.project, shouldUpdates)
	return
}

func updateAnyswapMatchTag(crossIns model.Datas) (shouldUpdates model.Datas, errs []*error) {
	var isStringAlphabetic = regexp.MustCompile(`^[0-9]+$`).MatchString
	// 若包含字母则返回false，不包含字母则返回true

	for _, crossIn := range crossIns {
		s := crossIn.MatchTag

		if ert := isStringAlphabetic(s[2:]); !ert { //是更新前的形式，即srcTxHash，需要进一步处理
			var swapIDHash common.Hash
			if utils.IsHex(s) {
				swapIDHash = common.HexToHash(s)
			} else {
				swapIDHash = common.BytesToHash([]byte(s))
			}
			crossIn.MatchTag = swapIDHash.String()
			shouldUpdates = append(shouldUpdates, crossIn)
		}
	}
	return
}

func (a *SimpleInMatcher) SendMail(sub string, datas []*model.Data) {
	subject := fmt.Sprintf("%s %s", strings.ToUpper(a.project), strings.ToUpper(sub))
	var info string
	for _, d := range datas {
		info += fmt.Sprintf("%s tx, Id: %d, chain: %s, hash: %s\n", subject, d.Id, d.Chain, d.Hash)
	}

	err := utils.SendMail(subject, info)
	if err != nil {
		e := fmt.Errorf(info)
		errs := []*error{&e}
		utils.LogError(errs, "./risk.log")
	}
}

func isDuplicated(b, c *model.Data) bool {
	if (b.Hash == c.Hash && b.Number != c.Number) ||
		(b.Hash == c.Hash && b.Number == c.Number && b.Index == c.Index) {
		return true
	}
	return false
}
