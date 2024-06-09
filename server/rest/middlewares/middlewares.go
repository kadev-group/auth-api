package middlewares

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/metrics"
	"github.com/doxanocap/pkg/ctxholder"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Middlewares struct {
	service interfaces.IService
	metrics *metrics.APIMetrics
	log     *zap.Logger
}

func InitMiddlewares(service interfaces.IService, metrics *metrics.APIMetrics, log *zap.Logger) *Middlewares {
	return &Middlewares{
		service: service,
		log:     log,
		metrics: metrics,
	}
}

func (m *Middlewares) VerifySession(c *gin.Context) {
	log := m.log.Named("[SESSION]")
	token := m.getAuthToken(c)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.HttpUnauthorized)
		return
	}

	uSession, err := m.service.Auth().ValidateAccessToken(c, token)
	if err != nil {
		httpError := errs.UnmarshalError(err)
		log.Error(err.Error())

		if httpError.StatusCode == 0 {
			c.AbortWithStatus(httpError.StatusCode)
			return
		}

		m.metrics.ErrorSessionVerification.Inc()
		c.AbortWithStatusJSON(httpError.StatusCode, err)
		return
	}

	m.metrics.SuccessSessionVerification.Inc()
	ctxholder.SetUserID(c, uSession.UserIDCode)
	c.Next()
}

func (m *Middlewares) GinMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

func (m *Middlewares) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)

		route := c.FullPath()
		status := c.Writer.Status()
		if route == "" && status == http.StatusNotFound {
			return
		}

		fields := []zap.Field{
			zap.String("route", c.FullPath()),
			zap.Duration("latency", latency),
			zap.Int("status", c.Writer.Status()),
		}

		log := m.log.With(zap.Any("payload", fields))

		privateErr := errs.GetGinPrivateErr(c)
		if privateErr != nil {
			m.metrics.ErrorHttpRequests.Inc()
			log.Error(privateErr.Error())
			if c.Writer.Status() == http.StatusOK {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		// public error will be set automatically if errs.SetGinError is called with it
		m.metrics.SuccessfulHttpRequests.Inc()
		log.Info("ok")
	}
}

func (m *Middlewares) SetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ctxholder.GetRefreshToken(c)
		c.SetCookie(consts.RefreshTokenKey,
			token,
			int(consts.RefreshTokenTTL),
			"/",
			"localhost",
			false,
			true)
	}

}

func (m *Middlewares) getAuthToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return ""
	}
	return split[2]
}
