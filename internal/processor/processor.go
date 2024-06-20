package processor

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/processor/apis"
	"auth-api/internal/processor/cache"
	"auth-api/internal/processor/queue"
	"auth-api/internal/processor/sms"
	"go.uber.org/zap"
)

type Processor struct {
	cacheProcessor interfaces.ICacheProcessor
	queueProcessor interfaces.IQueueProcessor
	apisProcessor  interfaces.IAPIsProcessor
	smsProcessor   interfaces.ISMSProcessor
}

func InitProcessor(
	log *zap.Logger,
	config *models.Config,
	smsProvider interfaces.ISMSProvider,
	cacheProvider interfaces.ICacheProvider,
	queueConsumerProvider interfaces.IQueueProducerProvider) *Processor {
	return &Processor{
		smsProcessor:   sms.NewSMSProcessor(smsProvider, log.Named("[SMS]")),
		queueProcessor: queue.NewQueueProcessor(config, queueConsumerProvider, log.Named("[QUEUE]")),
		apisProcessor:  apis.NewAPIsProcessor(config, log.Named("[APIs]")),
		cacheProcessor: cache.NewCacheProcessor(cacheProvider, log.Named("[CACHE]")),
	}
}

func (p *Processor) Cache() interfaces.ICacheProcessor {
	return p.cacheProcessor
}

func (p *Processor) Queue() interfaces.IQueueProcessor {
	return p.queueProcessor
}

func (p *Processor) SMS() interfaces.ISMSProcessor {
	return p.smsProcessor
}

func (p *Processor) APIs() interfaces.IAPIsProcessor {
	return p.apisProcessor
}
