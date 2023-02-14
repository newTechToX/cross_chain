package logic

import (
	"app/logic/aml"
	"app/logic/check_fake"
	"app/logic/replay"
	"app/model"
	"app/svc"
)

type Logic struct {
	svc      *svc.ServiceContext
	replayer *replay.Replayer
	aml      *aml.AML
	checker  *check_fake.Checker
}

func NewLogic(svc *svc.ServiceContext, chain string, config_path string) *Logic {
	c := check_fake.NewChecker(svc, chain, config_path)
	r := replay.NewReplayer(svc, c.Aml(), config_path)
	return &Logic{
		svc:      svc,
		replayer: r,
		checker:  c,
		aml:      c.Aml(),
	}
}

// fake token 和 fake chainId
//chainID的检查还没完成

func (a *Logic) CheckFake(project string, datas model.Datas) map[int]model.Datas {
	var fake = make(map[int]model.Datas)
	res := a.checker.IsFakeToken(project, datas)
	for _, d := range datas {
		if res[d.Token] != check_fake.SAFE {
			fake[res[d.Token]] = append(fake[res[d.Token]], d)
		}
	}
	return fake
}
