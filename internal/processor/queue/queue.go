package queue

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/processor/queue/producers"
	"go.uber.org/zap"
	"sync"
)

type Queue struct {
	log              *zap.Logger
	config           *models.Config
	consumerProvider interfaces.IQueueProducerProvider

	producersProcessor       interfaces.IQueueProducersProcessor
	producersProcessorRunner sync.Once
}

func NewQueueProcessor(
	config *models.Config,
	consumersProvider interfaces.IQueueProducerProvider,
	log *zap.Logger) *Queue {
	return &Queue{
		log:              log,
		config:           config,
		consumerProvider: consumersProvider,
	}
}

func (q *Queue) Producers() interfaces.IQueueProducersProcessor {
	q.producersProcessorRunner.Do(func() {
		q.producersProcessor = producers.NewProducersProcessor(q.config, q.consumerProvider, q.log.Named("[PRODUCERS]"))
	})
	return q.producersProcessor
}
