package logic

import (
	"app/logic/aml"
	"app/logic/check_fake"
	"app/logic/replay"
	"app/model"
	"app/svc"
)

const (
	interval  = 30 * 60
	batchSize = 10000
)

type Logic struct {
	svc      *svc.ServiceContext
	replayer *replay.Replayer
	aml      *aml.AML
	checker  *check_fake.Checker
}

func NewLogic(svc *svc.ServiceContext, chain string, config_path string) *Logic {
	if config_path == "" {
		config_path = "./txt_config.yaml"
	}
	c := check_fake.NewChecker(svc, chain, config_path)
	r := replay.NewReplayer(svc, c.Aml(), config_path)
	return &Logic{
		svc:      svc,
		replayer: r,
		checker:  c,
		aml:      c.Aml(),
	}
}

func (a *Logic) Main() {

}

func (a *Logic) CheckFake(project string, datas model.Datas) model.Datas {
	var fake_token = model.Datas{}
	res := a.checker.IsFakeToken(project, datas)
	for _, d := range datas {
		if res[d.Token] == check_fake.FAKE_TOKEN || res[d.Token] == check_fake.NULL_IN_DATABASE {
			fake_token = append(fake_token, d)
		}
	}
	return fake_token
}
