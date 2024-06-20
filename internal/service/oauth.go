package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/service/oauth"
)

type OAuthService struct {
	config  *models.Config
	manager interfaces.IManager

	gmail interfaces.IGmailService
}

func InitOAuthService(manager interfaces.IManager, config *models.Config) *OAuthService {
	return &OAuthService{
		config:  config,
		manager: manager,

		gmail: oauth.InitGmailService(manager, config),
	}
}

func (s *OAuthService) Gmail() interfaces.IGmailService {
	return s.gmail
}
