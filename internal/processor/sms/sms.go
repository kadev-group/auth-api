package sms

import (
	"auth-api/internal/interfaces"
	"context"
	"go.uber.org/zap"
)

type SMS struct {
	log      *zap.Logger
	provider interfaces.ISMSProvider
}

func NewSMSProcessor(provider interfaces.ISMSProvider, log *zap.Logger) *SMS {
	return &SMS{
		log:      log,
		provider: provider,
	}
}

func (sp *SMS) Send(ctx context.Context, phone, message string) error {
	log := sp.log.With(
		zap.String("phone_number", phone),
		zap.String("message", message))

	err := sp.provider.Send(ctx, phone, message)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
