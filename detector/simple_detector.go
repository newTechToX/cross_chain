package detector

import (
	"app/logic"
	"app/matcher"
	"app/model"
	"app/svc"
	"app/utils"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	log2 "log"
	"math/big"
)

type SimpleOutDetector struct {
	svc     *svc.ServiceContext
	logic   *logic.Logic
	matcher *matcher.Matcher
}

var _ model.Detector = &SimpleOutDetector{}

func NewSimpleOutDetector(svc *svc.ServiceContext, chain string) *SimpleOutDetector {
	path := "../logic/txt_config.yaml"
	return &SimpleOutDetector{
		svc:     svc,
		logic:   logic.NewLogic(svc, chain, path),
		matcher: matcher.NewMatcher(svc),
	}
}

// OutDetector的 Detect 用于检测所有tx的fake token & chainID，将没有match的做二次检测
//

func (a *SimpleOutDetector) Detect(project string, crossOuts []*model.Data) (riskTxs model.Datas, err error) {
	fake_token := a.logic.CheckFake(project, crossOuts)
	if len(fake_token) != 0 {
		for _, f := range fake_token {
			log2.Prefix()
			log.Warn("fake token detected", "project", project, "token", f.Token, "hash", f.Hash)
		}
	}

	for _, crossOut := range crossOuts {
		if crossOut.Direction != model.OutDirection {
			log.Warn("matching should not input cross-out")
			continue
		}

		var pending model.Datas
		stmt := fmt.Sprintf("select * from %s where direction = '%s' and isfaketoken != 1 and match_id is null", project, model.OutDirection)
		err := a.svc.Dao.DB().Select(&pending, stmt)
		if err != nil {
			return nil, err
		}
		if len(pending) == 0 {
			continue
		}
		// if len(pending) > 1 {
		// 	log.Error("multi matched", "src", crossIn.Hash)
		// }
		valid := make(model.Datas, 0)
		for _, counterparty := range pending {
			if !isMatched(counterparty, crossIn) {
				continue
			}
			valid = append(valid, counterparty)
			fillEmptyFields(counterparty, crossIn)
		}
		if len(valid) > 1 {
			log.Error("multi matched", "src", crossIn.Hash, "chain", crossIn.Chain, "project", project)
		} else {
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
	if out.FromAddress != "" && in.FromAddress != "" && out.FromAddress != in.FromAddress {
		return false
	}
	if out.ToAddress != "" && in.ToAddress != "" && out.ToAddress != in.ToAddress {
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
