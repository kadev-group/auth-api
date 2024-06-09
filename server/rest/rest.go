package rest

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"auth-api/server/rest/controllers"
	"auth-api/server/rest/middlewares"
	"context"
	"errors"
	"fmt"
	"github.com/doxanocap/pkg/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type REST struct {
	log     *zap.Logger
	config  *models.Config
	router  *gin.Engine
	server  *http.Server
	service interfaces.IService

	user        *controllers.UserController
	web         *controllers.WebController
	mobile      *controllers.MobileController
	session     *controllers.SessionController
	middlewares *middlewares.Middlewares
}

func InitREST(config *models.Config, service interfaces.IService, log *zap.Logger) *REST {
	m := metrics.NewAPIMetrics()
	return &REST{
		log:     log,
		config:  config,
		service: service,
		router:  router.InitGinRouter(config.ENV),

		user:    controllers.InitUserController(config, service, m, log.Named("[USER]")),
		web:     controllers.InitWebController(config, service, m, log.Named("[WEB]")),
		mobile:  controllers.InitMobileController(config, service, m, log.Named("[MOBILE]")),
		session: controllers.InitSessionController(config, service, m, log.Named("[SESSION]")),

		middlewares: middlewares.InitMiddlewares(service, m, log.Named("[MIDDLEWARE]")),
	}
}

func (r *REST) Run() {
	r.InitRoutes()
	r.server = &http.Server{
		Addr:           ":" + r.config.ServerPORT,
		Handler:        r.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		r.log.Info(fmt.Sprintf("REST server running at: %s", r.config.ServerPORT))
		if err := r.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			r.log.Error(fmt.Sprintf("r.ListenAndServer: %v", err))
		}
	}()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-ch

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := r.server.Shutdown(ctx); err != nil {
			r.log.Error(fmt.Sprintf("r.server.Stop: %s", err))
		}
		r.log.Info("REST graceful shut down...")
	}()
}
