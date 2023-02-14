package detector

import (
	"app/logic"
	"app/matcher"
	"app/model"
	"app/svc"
	"github.com/ethereum/go-ethereum/log"
	log2 "log"
)

type SimpleOutDetector struct {
	svc     *svc.ServiceContext
	logic   *logic.Logic
	matcher *matcher.Matcher
}

var _ model.Detector = &SimpleOutDetector{}

func NewSimpleOutDetector(svc *svc.ServiceContext, chain string, config_path string) *SimpleOutDetector {
	return &SimpleOutDetector{
		svc:     svc,
		logic:   logic.NewLogic(svc, chain, config_path),
		matcher: matcher.NewMatcher(svc),
	}
}

// OutDetector的 Detect 用于检测所有tx的fake token & chainID，将没有match的做二次检测
// 对于fake chainID的检查还没做完

func (a *SimpleOutDetector) DetectFake(project string, crossOuts []*model.Data) int {
	fake_token := a.logic.CheckFake(project, crossOuts)

	if len(fake_token) != 0 {
		for tag, f := range fake_token {
			for _, d := range f {
				log2.Prefix()
				log.Warn("fake token detected", "project", project, "tag", tag, "token", d.Token, "hash", d.Hash)
			}

		}
	}
	return len(fake_token)
}
