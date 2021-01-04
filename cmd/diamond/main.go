package main

import (
	"context"
	"fmt"
	"github.com/senayst4745/diamond/internal/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xlab/closer"
)

func main() {
	defer closer.Close()

	closer.Bind(func() {
		log.Info().Msg("shutdown")
	})

	cfg, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Can not init config")
	}

	if err := initLogger(cfg); err != nil {
		log.Fatal().Err(err).Msg("Can not init logger")
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	closer.Bind(cancelCtx)

	a, cleanup, err := initApp(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).
			Msg("Can't init app")
	}
	closer.Bind(cleanup)
	closer.Bind(func() {
		if err := a.Close(); err != nil {
			log.Error().Err(err).Msg("can not stop web application")
		}
	})
	if err := a.Start(); err != nil {
		log.Fatal().Err(err).Msg("Can not start app")
	}
}

func initLogger(c *config) error {
	log.Debug().Msg("initialize logger")
	logLvl, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLvl)
	switch c.LogFmt {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	case "json":
	default:
		return fmt.Errorf("unknown output format %s", c.LogFmt)
	}
	return nil
}

func initMongoConnection(ctx context.Context, cfg *config) (*mongo.Database, func(), error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DbAddr))
	if err != nil {
		return nil, nil, err
	}

	// Create connect
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return client.Database(cfg.DbName), func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Error().Err(err).Msg("error while connect to mongo")
		}
	}, nil
}

func initApiConfig(cfg *config) *api.Config {
	return &api.Config{Addr: cfg.Listen}
}
