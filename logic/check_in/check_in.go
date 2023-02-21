package check_in

import (
	"app/dao"
	"app/model"
	"fmt"
)

type InChecker struct {
	dao *dao.Dao
}

func NewInChecker(dao *dao.Dao) *InChecker {
	return &InChecker{
		dao,
	}
}

//规则1：检查是否有重复
//规则2：定时将unmatch的

func (a *InChecker) HasDuplicates(project string, datas model.Datas) map[uint64][]uint64 {
	var dup_map = make(map[uint64][]uint64)
	for _, data := range datas {
		stmt := fmt.Sprintf("select %s from %s where direction = 'in' and chain = '%s' and to_address = '%s' and from_chain = %s and to_chain = %s and id != %d and amount = %s",
			model.ResultRows, project, data.Chain, data.ToAddress, data.FromChainId.String(), data.ToChainId.String(), data.Id, data.Amount.String())
		var dups model.Datas
		err := a.dao.DB().Select(&dups, stmt)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, d := range dups {
				dup_map[data.Id] = append(dup_map[data.Id], d.Id)
			}
		}
	}
	return dup_map
}
