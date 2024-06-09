package producers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"go.uber.org/zap"
	"sync"
)

type Producers struct {
	consumerProvider interfaces.IQueueProducerProvider

	mailsProducer       interfaces.IQueueProducerProcessor
	mailsProducerRunner sync.Once
}

func NewProducersProcessor(
	config *models.Config,
	consumerProvider interfaces.IQueueProducerProvider,
	log *zap.Logger) *Producers {
	return &Producers{
		mailsProducer: NewMailsProducer(config, consumerProvider, log.Named("[MAILs]")),
	}
}

func (p *Producers) Mails() interfaces.IQueueProducerProcessor {
	return p.mailsProducer
}
