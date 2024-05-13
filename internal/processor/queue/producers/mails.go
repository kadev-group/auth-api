package producers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"context"
	"encoding/json"
	"go.uber.org/zap"
)

type MailsProducer struct {
	log      *zap.Logger
	config   *models.Config
	provider interfaces.IQueueProducerProvider
}

func NewMailsProducer(
	config *models.Config,
	provider interfaces.IQueueProducerProvider,
	log *zap.Logger) *MailsProducer {
	return &MailsProducer{
		log:      log,
		config:   config,
		provider: provider,
	}
}

func (mc *MailsProducer) Send(ctx context.Context, message *models.MailsProducerMsg) error {
	log := mc.log.With(
		zap.String("send_to", message.SendTo),
		zap.String("code", message.VerificationCode)).
		Named("Send")

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err = mc.provider.Send(ctx, mc.config.MailsQueue, body); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
