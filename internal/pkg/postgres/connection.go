package postgres

import (
	"auth-api/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func getDSN(cfg models.PSQL) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PqHOST, cfg.PqPORT, cfg.PqUSER, cfg.PqPASSWORD, cfg.PqDATABASE, cfg.PqSSL)
}

func InitConnection(config *models.Config, log *zap.Logger) *sqlx.DB {
	log = log.Named("[PSQL]")
	conn, err := sqlx.Connect("postgres", getDSN(config.PSQL))
	if err != nil {
		log.Fatal(fmt.Sprintf("connect: %s", err))
	}

	if err = conn.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("ping: %s", err))
	}

	return conn
}
