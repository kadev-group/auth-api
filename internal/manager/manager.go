package manager

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/rabbitmq"
	"auth-api/internal/pkg/redis"
	"auth-api/internal/pkg/smsc"
	"auth-api/internal/processor"
	"auth-api/internal/repository"
	"auth-api/internal/service"
	"auth-api/server"
	_ "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Manager struct {
	db    *sqlx.DB
	cache *redis.Conn

	service    interfaces.IService
	repository interfaces.IRepository
	processor  interfaces.IProcessor
	server     interfaces.IServer
}

func InitManager(
	db *sqlx.DB,
	log *zap.Logger,
	cache *redis.Conn,
	config *models.Config,
	smsProvider *smsc.SMSc,
	queueProducer *rabbitmq.ProducerClient) *Manager {
	m := &Manager{
		db:    db,
		cache: cache,
	}

	m.service = service.InitService(m, config)
	m.processor = processor.InitProcessor(log, config, smsProvider, cache, queueProducer)
	m.repository = repository.InitRepository(db, config, m.Processor().Cache())
	m.server = server.InitServer(log, config, m.Service())

	return m
}

func (m *Manager) Repository() interfaces.IRepository {
	return m.repository
}

func (m *Manager) Service() interfaces.IService {
	return m.service
}

func (m *Manager) Processor() interfaces.IProcessor {
	return m.processor
}

func (m *Manager) Server() interfaces.IServer {
	return m.server
}
