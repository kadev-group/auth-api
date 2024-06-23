package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/metrics"
	"fmt"
	"github.com/doxanocap/pkg/ctxholder"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type SessionController struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService
}

func InitSessionController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *SessionController {
	return &SessionController{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

func (ctl *SessionController) Refresh(c *gin.Context) {
	ctl.metrics.RefreshRequest.Inc()

	token, err := c.Cookie(consts.RefreshTokenKey)
	if err != nil {
		errs.SetBothErrors(c, models.HttpUnauthorized, err)
		return
	}

	ctxholder.SetClientIP(c)
	response, err := ctl.service.User().Refresh(c, token)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctxholder.SetRefreshToken(c, response.RefreshToken, 0)
	c.JSON(http.StatusOK, response)
}

func (ctl *SessionController) Logout(c *gin.Context) {
	ctl.metrics.LogoutRequest.Inc()

	token, err := c.Cookie(consts.RefreshTokenKey)
	if err != nil {
		errs.SetBothErrors(c, models.HttpUnauthorized, err)
		return
	}

	if err = ctl.service.User().Logout(c, token); err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctxholder.SetRefreshToken(c, "", -1)
	c.Status(http.StatusOK)
}

func (ctl *SessionController) Verify(c *gin.Context) {
	ctl.metrics.VerifySessionRequest.Inc()
	log := ctl.log.Named("[VERIFY]")

	token := ctl.getAccessToken(c)
	if token == "" {
		errs.SetGinError(c, models.HttpUnauthorized)
		return
	}

	userSession, err := ctl.service.Auth().ValidateAccessToken(c, token)
	if err != nil {
		httpError := errs.UnmarshalError(err)
		log.Error(fmt.Sprintf("ValidateAccessToken: %s", err))

		switch httpError.StatusCode {
		case http.StatusInternalServerError:
			c.Status(http.StatusUnauthorized)
		default:
			c.JSON(http.StatusUnauthorized, err)
		}
		return
	}

	c.JSON(http.StatusOK, userSession)
}

func (ctl *SessionController) SendValidateCode(c *gin.Context) {
	ctl.metrics.VerifyEmailRequests.Inc()

	var request models.SendValidateCodeReq
	if err := c.ShouldBindJSON(&request); err != nil {
		errs.SetBothErrors(c, models.HttpBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		errs.SetGinError(c, err)
		return
	}

	err := ctl.service.User().SendValidateCode(c, &request)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (ctl *SessionController) getAccessToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}
