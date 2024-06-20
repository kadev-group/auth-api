package apis

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"go.uber.org/zap"
)

type APIs struct {
	googleAPI interfaces.IGoogleAPIProcessor
}

func NewAPIsProcessor(config *models.Config, log *zap.Logger) *APIs {
	return &APIs{
		googleAPI: InitGoogleAPI(config, log.Named("[AUTH_API]")),
	}
}

func (a *APIs) GoogleAPI() interfaces.IGoogleAPIProcessor {
	return a.googleAPI
}
