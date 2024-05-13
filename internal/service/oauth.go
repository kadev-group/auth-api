package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/service/oauth"
)

type OAuthService struct {
	config  *models.Config
	manager interfaces.IManager

	google interfaces.IGoogleAPI
}

func InitOAuthService(manager interfaces.IManager, config *models.Config) *OAuthService {
	return &OAuthService{
		config:  config,
		manager: manager,

		google: oauth.InitGoogleAPI(manager, config),
	}
}

func (s *OAuthService) Google() interfaces.IGoogleAPI {
	return s.google
}
