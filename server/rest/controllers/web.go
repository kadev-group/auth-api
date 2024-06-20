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

const (
	keyRefreshToken = "refresh_token"
)

type WebController struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService
}

func InitWebController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *WebController {
	return &WebController{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

func (ctl *WebController) SignIn(c *gin.Context) {
	ctl.metrics.SignInRequests.Inc()

	var request models.WebSignInReq
	if err := c.ShouldBindJSON(&request); err != nil {
		errs.SetBothErrors(c, models.HttpBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctxholder.SetClientIP(c)
	response, err := ctl.service.User().Authenticate(c, request.ToAuthReq())
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctl *WebController) SignUp(c *gin.Context) {
	ctl.metrics.SignUpRequests.Inc()

	var request models.SignUpReq
	if err := c.ShouldBindJSON(&request); err != nil {
		errs.SetBothErrors(c, models.HttpBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctxholder.SetClientIP(c)
	response, err := ctl.service.User().Create(c, request.ToUserDTO())
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GmailAuthRedirect ...
func (ctl *WebController) GmailAuthRedirect(c *gin.Context) {
	ctl.metrics.GoogleRedirectRequest.Inc()
	log := ctl.log.Named("GoogleRedirect")
	state := c.Param("state")
	if !tools.IsUUID(state) {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	response, err := ctl.service.OAuth().Gmail().GetRedirectURL(c, state)
	if err != nil {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	log.Info(fmt.Sprintf("redirected | %s", state))
	c.JSON(http.StatusOK, response)
}

// GmailAuthCallBack ...
func (ctl *WebController) GmailAuthCallBack(c *gin.Context) {
	ctl.metrics.GoogleCallBackRequest.Inc()
	log := ctl.log.Named("GoogleCallBack")

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

	//ctxholder.SetKV(c, "client_ip", c.ClientIP())
	accessToken, err := ctl.service.OAuth().Gmail().HandleCallBack(c, state, exchangeCode)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	redirectURL, err := url.Parse(ctl.config.OAuth.ClientCallBackURI)
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
