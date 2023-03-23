package check_fake

import (
	"app/logic/aml"
	"app/model"
	"app/provider"
	"app/svc"
	"app/utils"
	"database/sql"
	"fmt"
	log2 "github.com/ethereum/go-ethereum/log"
	"log"
	"strings"
)

type FakeChecker struct {
	svc      *svc.ServiceContext
	aml      *aml.AML
	provider *provider.Provider
}

const (
	NULL_IN_DATABASE = -1
	FAKE_TOKEN       = 1
	SAFE             = 0
	FAKE_CHAINID     = 2
)

func NewFakeChecker(svc *svc.ServiceContext, chain string, config_path string) *FakeChecker {
	p := svc.Providers.Get(chain)
	if p == nil {
		panic(fmt.Sprintf("%v: invalid provider", chain))
	}
	return &FakeChecker{
		svc:      svc,
		aml:      aml.NewAML(config_path),
		provider: p,
	}
}

func (c *FakeChecker) Aml() *aml.AML {
	if c.aml == nil {
		c.aml = aml.NewAML("../txt_config.yaml")
	}
	return c.aml
}

func (a *FakeChecker) IsFakeToken(project string, t *model.Data) int {
	res := a.queryTable(project, t)
	var errs []*error

	switch res {
	//如果不在数据库中，就查询aml
	case NULL_IN_DATABASE:
		tw := []string{t.Token}
		info_from_aml, err := a.aml.QueryAml(t.Chain, tw)
		if err != nil {
			s := fmt.Sprintf("failed to query token from aml, chain:%s, address:%s", t.Chain, t.Token)
			log.SetPrefix("IsFakeToken()")
			utils.LogPrint(s, "./logic.log")
		}

		//如果aml里面也查不到token，就查deployers
		if info_from_aml[t.Token] == nil {
			log2.Warn("IsFakeToken(), failed to query token from aml ", "Project", project, "Chain", t.Chain, "Token", t.Token)
			a.provider = a.svc.Providers.Get(t.Chain)
			deployer_info_from_provider, err := a.provider.GetContractInfo(t.Token)
			if err != nil || deployer_info_from_provider == nil {
				s := fmt.Sprintf("failed to query token from providers, chain:%s, address:%s", t.Chain, t.Token)
				log.SetPrefix("IsFakeToken()")
				utils.LogPrint(s, "./logic.log")
				break
			}
			tw = []string{deployer_info_from_provider.Deployer}
			deployer_info_aml, err := a.aml.QueryAml(t.Chain, tw)
			if err != nil {
				s := fmt.Sprintf("failed to query token from providers, chain:%s, address:%s", t.Chain, t.Token)
				log.SetPrefix("IsFakeToken()")
				utils.LogPrint(s, "./logic.log")
				break
			}

			//如果aml也查不到deployer的信息，则认为风险较高
			if deployer_info_aml[deployer_info_from_provider.Deployer] == nil {
				s := fmt.Errorf("nothing: query deployer from providers, chain:%s, address:%s", t.Chain, deployer_info_from_provider.Deployer)
				log.SetPrefix("IsFakeToken()")
				errs = append(errs, &s)
			} else { // 如果查到了deployer的信息，若name前缀 == "Multichain"
				if deployer_info_aml[deployer_info_from_provider.Deployer][0].Name[:10] == "Multichain" || deployer_info_aml[deployer_info_from_provider.Deployer][0].Name == "CONTRACT DEPLOYER" {
					res = SAFE
				} else {
					res = FAKE_TOKEN
					s := fmt.Errorf("IsFakeToken(): deployer risk, chain:%s, address:%s, name:%s",
						t.Chain, deployer_info_from_provider.Deployer, deployer_info_aml[deployer_info_from_provider.Deployer][0].Name)
					errs = append(errs, &s)
				}
			}
		} else if info_from_aml[t.Token][0].Risk >= 3 { // 如果从aml里查到了token的标签
			res = FAKE_TOKEN
			s := fmt.Errorf("IsFakeToken(): token risk %d, chain:%s, address:%s, name:%s",
				info_from_aml[t.Token][0].Risk, t.Chain, t.Token, info_from_aml[t.Token][0].Name)
			errs = append(errs, &s)
		} else if project == "across" && len(info_from_aml[t.Token][0].Labels) > 0 {
			for _, l := range info_from_aml[t.Token][0].Labels {
				if strings.Contains(l, "TOKEN") {
					stmt := fmt.Sprintf("update %s set isfaketoken=2 where id=%d", project, t.Id)
					_, err = a.svc.Dao.DB().Exec(stmt)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else {
			res = SAFE
		}
		break

	case FAKE_TOKEN:
		s := fmt.Errorf("fake token: fake token in database, chain:%s, address:%s, hash:%s", t.Chain, t.Token, t.Hash)
		errs = append(errs, &s)
		break

	case SAFE:
		break
	}
	utils.LogError(errs, "./risk.log")
	return res
}

func (a *FakeChecker) queryTable(project string, t *model.Data) int {
	stmt := fmt.Sprintf("select isfaketoken from %s where chain = '%s' and token = '%s' and block_number<%d", project, t.Chain, t.Token, t.Number)
	var fake sql.NullInt32
	err := a.svc.ProjectsDao.DB().Get(&fake, stmt)
	if err != nil {
		s := fmt.Sprintf("Check_token(): failed to query token from table, chain:%s, address:%s", t.Chain, t.Token)
		utils.LogPrint(s, "./logic.log")
		return NULL_IN_DATABASE
	}

	if fake.Valid && fake.Int32 == 1 {
		return FAKE_TOKEN
	} else {
		return SAFE
	}
}
