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

type Check struct {
	svc      *svc.ServiceContext
	aml      *aml.AML
	provider *provider.Provider
}

func NewChecker(svc *svc.ServiceContext, chain, path string) *Check {
	p := svc.Providers.Get(chain)
	if p == nil {
		panic(fmt.Sprintf("%v: invalid provider", chain))
	}
	return &Check{
		svc:      svc,
		aml:      aml.NewAML(path),
		provider: p,
	}
}

func (c *Check) Aml() *aml.AML {
	if c.aml == nil {
		c.aml = aml.NewAML("../txt_config.yaml")
	}
	return c.aml
}

const null_in_database, is_fake, not_fake = -1, 1, 0

func (a *Check) IsFakeToken(project string, tokens model.TokenChains) map[string]int {
	var res = make(map[string]int)
	for _, t := range tokens {
		res_table := a.queryTable(project, t)
		res[t.Token] = res_table
		a.provider = (&provider.Providers{}).Get(t.Chain)

		switch res_table {
		//如果不在数据库中，就查询aml
		case null_in_database:
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
						res[t.Token] = not_fake
					} else {
						res[t.Token] = is_fake
						s := fmt.Sprintf("IsFakeToken(): deployer risk, chain:%s, address:%s, name:%s",
							t.Chain, deployer_info_from_provider.Deployer, deployer_info_aml[deployer_info_from_provider.Deployer][0].Name)
						log.Println(s)
						utils.LogPrint(s, "./risk.log")
					}
				}
			} else if info_from_aml[t.Token][0].Risk >= 3 { // 如果从aml里查到了token的标签
				res[t.Token] = is_fake
				s := fmt.Sprintf("IsFakeToken(): token risk %d, chain:%s, address:%s, name:%s",
					info_from_aml[t.Token][0].Risk, t.Chain, t.Token, info_from_aml[t.Token][0].Name)
				log.Println(s)
				utils.LogPrint(s, "./risk.log")
			} else {
				res[t.Token] = not_fake
			}
			break

		case is_fake:
			s := fmt.Sprintf("IsFakeToken(): fake token: fake token in database, chain:%s, address:%s", t.Chain, t.Token)
			log.Println(s)
			utils.LogPrint(s, "./risk.log")
			break

		case not_fake:
			break
		}
	}
	return res
}

func (a *Check) queryTable(project string, t *model.TokenChain) int {
	stmt := fmt.Sprintf("select isfaketoken from %s where chain = '%s' and token = '%s' and block_number<%d", project, t.Chain, t.Token, t.Block)
	var fake sql.NullInt32
	err := a.svc.ProjectsDao.DB().Get(&fake, stmt)
	if err != nil {
		s := fmt.Sprintf("Check_token(): failed to query token from table, chain:%s, address:%s", t.Chain, t.Token)
		utils.LogPrint(s, "./logic.log")
		return null_in_database
	}

	if fake.Valid && fake.Int32 == 1 {
		return is_fake
	} else {
		return not_fake
	}
}
