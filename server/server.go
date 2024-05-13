package server

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/server/rest"
	"go.uber.org/zap"
	"sync"
)

type Server struct {
	log     *zap.Logger
	config  *models.Config
	service interfaces.IService

	restServer       interfaces.IRESTServer
	restServerRunner sync.Once
}

func InitServer(
	log *zap.Logger,
	config *models.Config,
	service interfaces.IService) *Server {
	return &Server{
		log:     log,
		config:  config,
		service: service,
	}
}

func (p *Server) REST() interfaces.IRESTServer {
	p.restServerRunner.Do(func() {
		p.restServer = rest.InitREST(p.config, p.service, p.log.Named("[REST]"))
	})
	return p.restServer
}
