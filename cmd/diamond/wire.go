//+build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/senayst4745/diamond/internal/api"
	"github.com/senayst4745/diamond/internal/repository"
	"github.com/senayst4745/diamond/internal/service"
)

func initApp(ctx context.Context, cfg *config) (a *api.API, closer func(), err error) {
	wire.Build(initApiConfig, initMongoConnection,
		repository.New,
		wire.Bind(new(repository.MineRepository), new(*repository.MongoMineRepository)),
		service.New,
		wire.Bind(new(service.Service), new(*service.SimpleMineService)),
		api.New)
	return nil, nil, nil
}
