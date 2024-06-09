package queue

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/processor/queue/producers"
	"go.uber.org/zap"
	"sync"
)

type Queue struct {
	producersProcessor       interfaces.IQueueProducersProcessor
	producersProcessorRunner sync.Once
}

func NewQueueProcessor(
	config *models.Config,
	consumerProvider interfaces.IQueueProducerProvider,
	log *zap.Logger) *Queue {
	return &Queue{
		producersProcessor: producers.NewProducersProcessor(config, consumerProvider, log.Named("[PRODUCERS]")),
	}
}

func (q *Queue) Producers() interfaces.IQueueProducersProcessor {
	return q.producersProcessor
}
