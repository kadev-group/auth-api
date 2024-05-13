package producers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"go.uber.org/zap"
	"sync"
)

type Producers struct {
	log              *zap.Logger
	config           *models.Config
	consumerProvider interfaces.IQueueProducerProvider

	mailsProducer       interfaces.IQueueProducerProcessor
	mailsProducerRunner sync.Once
}

func NewProducersProcessor(
	config *models.Config,
	consumersProvider interfaces.IQueueProducerProvider,
	log *zap.Logger) *Producers {
	return &Producers{
		log:              log,
		config:           config,
		consumerProvider: consumersProvider,
	}
}

func (p *Producers) Mails() interfaces.IQueueProducerProcessor {
	p.mailsProducerRunner.Do(func() {
		p.mailsProducer = NewMailsProducer(p.config, p.consumerProvider, p.log.Named("[MAILs]"))
	})
	return p.mailsProducer
}
