package postgres

import (
	"auth-api/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func sslDisabled(config string) string {
	return fmt.Sprintf("%s?sslmode=disable", config)
}

func InitConnection(config *models.Config, log *zap.Logger) *sqlx.DB {
	log = log.Named("[PSQL]")
	conn, err := sqlx.Connect("postgres", sslDisabled(config.PsqlDsn))
	if err != nil {
		log.Fatal(fmt.Sprintf("connect: %s", err))
	}

	if err = conn.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("ping: %s", err))
	}

	return conn
}
