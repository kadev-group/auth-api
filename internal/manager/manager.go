package manager

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/rabbitmq"
	"auth-api/internal/pkg/redis"
	"auth-api/internal/processor"
	"auth-api/internal/repository"
	"auth-api/internal/service"
	"auth-api/server"
	_ "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	db                    *sqlx.DB
	log                   *zap.Logger
	config                *models.Config
	cacheProvider         interfaces.ICacheProvider
	queueProducerProvider interfaces.IQueueProducerProvider

	service       interfaces.IService
	serviceRunner sync.Once

	repository       interfaces.IRepository
	repositoryRunner sync.Once

	processor       interfaces.IProcessor
	processorRunner sync.Once

	server       interfaces.IServer
	serverRunner sync.Once
}

func InitManager(
	db *sqlx.DB,
	log *zap.Logger,
	config *models.Config,
	rabbitmqProducer *rabbitmq.ProducerClient,
	redisConn *redis.Conn) *Manager {
	return &Manager{
		db:                    db,
		log:                   log,
		config:                config,
		cacheProvider:         redisConn,
		queueProducerProvider: rabbitmqProducer,
	}
}

func (m *Manager) Repository() interfaces.IRepository {
	m.repositoryRunner.Do(func() {
		m.repository = repository.InitRepository(m.db, m.config, m.Processor().Cache())
	})
	return m.repository
}

func (m *Manager) Service() interfaces.IService {
	m.serviceRunner.Do(func() {
		m.service = service.InitService(m, m.config)
	})
	return m.service
}

func (m *Manager) Processor() interfaces.IProcessor {
	m.processorRunner.Do(func() {
		m.processor = processor.InitProcessor(m.log, m.config,
			m.Service(), m.queueProducerProvider, m.cacheProvider)
	})
	return m.processor
}

func (m *Manager) Server() interfaces.IServer {
	m.serverRunner.Do(func() {
		m.server = server.InitServer(m.log, m.config, m.Service())
	})
	return m.server
}

func (m *Manager) SetDB(db *sqlx.DB) {
	m.db = db
}

func (m *Manager) SetCache(cache interfaces.ICacheProvider) {
	m.cacheProvider = cache
}
