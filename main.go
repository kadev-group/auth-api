package main

import (
	"auth-api/internal/manager"
	"auth-api/internal/models"
	"auth-api/internal/pkg/postgres"
	"auth-api/internal/pkg/rabbitmq"
	"auth-api/internal/pkg/redis"
	"auth-api/internal/pkg/smsc"
	"github.com/doxanocap/pkg/config"
	"github.com/doxanocap/pkg/logger"
	"go.uber.org/fx"
	"log"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.InitConfig[models.Config],
			logger.InitLogger[models.Config],
			rabbitmq.NewProducerClient,
			smsc.NewSMSc,
			postgres.InitConnection,
			redis.InitConnection,
			manager.InitManager,
		),
		fx.Invoke(manager.Run),
	)

	app.Run()
	if err := app.Err(); err != nil {
		log.Fatal(err)
	}
}
