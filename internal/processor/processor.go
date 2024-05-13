package processor

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/processor/cache"
	"auth-api/internal/processor/queue"
	"go.uber.org/zap"
	"sync"
)

type Processor struct {
	log                   *zap.Logger
	config                *models.Config
	service               interfaces.IService
	cacheProvider         interfaces.ICacheProvider
	queueConsumerProvider interfaces.IQueueProducerProvider

	cacheProcessor       interfaces.ICacheProcessor
	cacheProcessorRunner sync.Once

	queueProcessor       interfaces.IQueueProcessor
	queueProcessorRunner sync.Once
}

func InitProcessor(
	log *zap.Logger,
	config *models.Config,
	service interfaces.IService,
	queueConsumerProvider interfaces.IQueueProducerProvider,
	cache interfaces.ICacheProvider) *Processor {
	return &Processor{
		log:                   log,
		cacheProvider:         cache,
		config:                config,
		service:               service,
		queueConsumerProvider: queueConsumerProvider,
	}
}

func (p *Processor) Cache() interfaces.ICacheProcessor {
	p.cacheProcessorRunner.Do(func() {
		p.cacheProcessor = cache.NewCacheProcessor(p.cacheProvider, p.log.Named("[CACHE]"))
	})
	return p.cacheProcessor
}

func (p *Processor) Queue() interfaces.IQueueProcessor {
	p.queueProcessorRunner.Do(func() {
		p.queueProcessor = queue.NewQueueProcessor(p.config, p.queueConsumerProvider, p.log.Named("[QUEUE]"))
	})
	return p.queueProcessor
}
