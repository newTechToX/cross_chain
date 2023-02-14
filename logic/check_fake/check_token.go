package check_fake

import (
	"app/logic/aml"
	"app/model"
	"app/provider"
	"app/svc"
	"app/utils"
	"database/sql"
	"fmt"
	"log"
)

type Checker struct {
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

func NewChecker(svc *svc.ServiceContext, chain string, config_path string) *Checker {
	p := svc.Providers.Get(chain)
	if p == nil {
		panic(fmt.Sprintf("%v: invalid provider", chain))
	}
	return &Checker{
		svc:      svc,
		aml:      aml.NewAML(config_path),
		provider: p,
	}
}

func (c *Checker) Aml() *aml.AML {
	if c.aml == nil {
		c.aml = aml.NewAML("../txt_config.yaml")
	}
	return c.aml
}

func (a *Checker) IsFakeToken(project string, tokens model.Datas) map[string]int {
	var res = make(map[string]int)
	for _, t := range tokens {
		res_table := a.queryTable(project, t)
		res[t.Token] = res_table
		a.provider = (&provider.Providers{}).Get(t.Chain)

		switch res_table {
		//如果不在数据库中，就查询aml
		case NULL_IN_DATABASE:
			tw := []string{t.Token}
			info_from_aml, err := a.aml.QueryAml(t.Chain, tw)
			if err != nil {
				s := fmt.Sprintf("IsFakeToken(): failed to query token from aml, chain:%s, address:%s", t.Chain, t.Token)
				utils.LogPrint(s, "./logic.log")
			}

			//如果aml里面也查不到token，就查deployers
			if info_from_aml[t.Token] == nil {
				s := fmt.Sprintf("IsFakeToken(): nothing: query token from providers, chain:%s, address:%s", t.Chain, t.Token)
				log.Println(s)
				deployer_info_from_provider, err := a.provider.GetContractInfo(t.Token)
				if err != nil {
					s := fmt.Sprintf("IsFakeToken(): failed to query token from providers, chain:%s, address:%s", t.Chain, t.Token)
					utils.LogPrint(s, "./logic.log")
					break
				}

				tw = []string{deployer_info_from_provider.Deployer}
				deployer_info_aml, err := a.aml.QueryAml(t.Chain, tw)
				if err != nil {
					s := fmt.Sprintf("IsFakeToken(): failed to query token from providers, chain:%s, address:%s", t.Chain, t.Token)
					utils.LogPrint(s, "./logic.log")
					break
				}

				//如果aml也查不到deployer的信息，则认为风险较高
				if deployer_info_aml[deployer_info_from_provider.Deployer] == nil {
					s := fmt.Sprintf("IsFakeToken(): nothing: query deployer from providers, chain:%s, address:%s", t.Chain, deployer_info_from_provider.Deployer)
					log.Println(s)
					utils.LogPrint(s, "./risk.log")
				} else { // 如果查到了deployer的信息，若name前缀 == "Multichain"
					if deployer_info_aml[deployer_info_from_provider.Deployer][0].Name[:10] == "Multichain" {
						res[t.Token] = SAFE
					} else {
						res[t.Token] = FAKE_TOKEN
						s := fmt.Sprintf("IsFakeToken(): deployer risk, chain:%s, address:%s, name:%s",
							t.Chain, deployer_info_from_provider.Deployer, deployer_info_aml[deployer_info_from_provider.Deployer][0].Name)
						log.Println(s)
						utils.LogPrint(s, "./risk.log")
					}
				}
			} else if info_from_aml[t.Token][0].Risk >= 3 { // 如果从aml里查到了token的标签
				res[t.Token] = FAKE_TOKEN
				s := fmt.Sprintf("IsFakeToken(): token risk %d, chain:%s, address:%s, name:%s",
					info_from_aml[t.Token][0].Risk, t.Chain, t.Token, info_from_aml[t.Token][0].Name)
				log.Println(s)
				utils.LogPrint(s, "./risk.log")
			} else {
				res[t.Token] = SAFE
			}
			break

		case FAKE_TOKEN:
			s := fmt.Sprintf("IsFakeToken(): fake token: fake token in database, chain:%s, address:%s", t.Chain, t.Token)
			log.Println(s)
			utils.LogPrint(s, "./risk.log")
			break

		case SAFE:
			println("safe")
			break
		}
	}
	return res
}

func (a *Checker) queryTable(project string, t *model.Data) int {
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
