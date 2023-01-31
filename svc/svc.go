package svc

import (
	"app/config"
	"app/dao"
	"app/provider"
	"context"
	"sync"
)

type ServiceContext struct {
	Ctx       context.Context
	Wg        sync.WaitGroup
	Config    *config.Config
	Providers *provider.Providers
	Dao       *dao.Dao
	//LabelDao    *dao.Dao
	ProjectsDao *dao.Dao
}

func NewServiceContext(ctx context.Context, c *config.Config) *ServiceContext {
	return &ServiceContext{
		Ctx:       ctx,
		Config:    c,
		Providers: provider.NewProviders(c),
		Dao:       dao.NewDao(c.Database.CrosschainDataSource),
		//LabelDao:    dao.NewLabelDao(c.Database.AddressLabelSource),
		ProjectsDao: dao.NewAnyDao(c.Database.AnyswapSource),
	}
}
