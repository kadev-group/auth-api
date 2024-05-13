package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"sync"
)

type Service struct {
	config  *models.Config
	manager interfaces.IManager

	auth       interfaces.IAuthService
	authRunner sync.Once

	oAuth       interfaces.IOAuthService
	oAuthRunner sync.Once

	user       interfaces.IUserService
	userRunner sync.Once
}

func InitService(manager interfaces.IManager, config *models.Config) *Service {
	return &Service{
		config:  config,
		manager: manager,
	}
}

func (s *Service) Auth() interfaces.IAuthService {
	s.authRunner.Do(func() {
		s.auth = InitAuthService(s.manager, s.config)
	})
	return s.auth
}

func (s *Service) User() interfaces.IUserService {
	s.userRunner.Do(func() {
		s.user = InitUserService(s.manager)
	})
	return s.user
}

func (s *Service) OAuth() interfaces.IOAuthService {
	s.oAuthRunner.Do(func() {
		s.oAuth = InitOAuthService(s.manager, s.config)
	})
	return s.oAuth
}
