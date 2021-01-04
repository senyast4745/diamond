// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"github.com/senayst4745/diamond/internal/api"
	"github.com/senayst4745/diamond/internal/repository"
	"github.com/senayst4745/diamond/internal/service"
)

// Injectors from wire.go:

func initApp(ctx context.Context, cfg *config) (*api.API, func(), error) {
	apiConfig := initApiConfig(cfg)
	database, cleanup, err := initMongoConnection(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	mongoMineRepository := repository.New(database)
	simpleMineService := service.New(mongoMineRepository)
	apiAPI, err := api.New(ctx, apiConfig, simpleMineService)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return apiAPI, func() {
		cleanup()
	}, nil
}