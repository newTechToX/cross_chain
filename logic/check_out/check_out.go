package replay

import "app/dao"

type OutChecker struct {
	dao *dao.Dao
}

func NewOutChecker(dao *dao.Dao) *OutChecker {
	return &OutChecker{
		dao,
	}
}
