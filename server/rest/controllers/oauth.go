package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"auth-api/internal/pkg/tools"
	"fmt"
	"github.com/doxanocap/pkg/ctxholder"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type OAuthController struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService
}

func InitOAuthController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *OAuthController {
	return &OAuthController{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

// GoogleRedirect ...
func (h *OAuthController) GoogleRedirect(c *gin.Context) {
	h.metrics.GoogleRedirectRequest.Inc()
	log := h.log.Named("GoogleRedirect")
	state := c.Param("state")
	if !tools.IsUUID(state) {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	response, err := h.service.OAuth().Google().GetRedirectURL(c, state)
	if err != nil {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	log.Info(fmt.Sprintf("redirected | %s", state))
	c.JSON(http.StatusOK, response)
}

// GoogleCallBack ...
func (h *OAuthController) GoogleCallBack(c *gin.Context) {
	h.metrics.GoogleCallBackRequest.Inc()
	log := h.log.Named("GoogleCallBack")

	errReason := c.Query("error_reason")
	if errReason != "" {
		errs.SetBothErrors(c, models.HttpBadRequest, errs.New(errReason))
		return
	}

	state := c.Query("state")
	exchangeCode := c.Query("code")
	if !tools.IsUUID(state) || exchangeCode == "" {
		errs.SetBothErrors(c, models.HttpBadRequest, errs.New("bad state or code"))
		return
	}

	ctxholder.SetKV(c, "client_ip", c.ClientIP())
	accessToken, err := h.service.OAuth().Google().HandleCallBack(c, state, exchangeCode)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	redirectURL, err := url.Parse(h.config.ClientCallBackURI)
	if err != nil {
		errs.SetBothErrors(c, models.HttpInternalServerError, err)
		return
	}
	query := redirectURL.Query()
	query.Add("state", state)
	query.Add("token", accessToken)
	redirectURL.RawQuery = query.Encode()

	log.Info(fmt.Sprintf("redirected | %s", query.Encode()))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}
