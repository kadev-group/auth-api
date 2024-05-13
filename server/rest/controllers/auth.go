package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/metrics"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
)

const (
	keyRefreshToken = "refresh_token"
)

type AuthController struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService
}

func InitAuthController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *AuthController {
	return &AuthController{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

func (ctl *AuthController) SignIn(c *gin.Context) {
	ctl.metrics.SignInRequests.Inc()

	var request models.SignInReq
	if err := c.ShouldBindJSON(&request); err != nil {
		errs.SetBothErrors(c, models.HttpBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		errs.SetGinError(c, err)
		return
	}

	response, err := ctl.service.User().Authenticate(c, request.ToUserDTO())
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctl.setRefreshToken(c, response.Tokens.RefreshToken)
	c.JSON(http.StatusOK, response)
}

func (ctl *AuthController) SignUp(c *gin.Context) {
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

	userDTO := request.ToUserDTO()
	response, err := ctl.service.User().Create(c, userDTO)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctl.setRefreshToken(c, response.Tokens.RefreshToken)
	c.JSON(http.StatusOK, response)
}

func (ctl *AuthController) Refresh(c *gin.Context) {
	ctl.metrics.RefreshRequest.Inc()

	token, err := c.Cookie(keyRefreshToken)
	if err != nil {
		errs.SetBothErrors(c, models.HttpUnauthorized, err)
		return
	}

	response, err := ctl.service.User().Refresh(c, token)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	fmt.Println(response)
	ctl.setRefreshToken(c, response.RefreshToken)
	c.JSON(http.StatusOK, response)
}

func (ctl *AuthController) Logout(c *gin.Context) {
	ctl.metrics.LogoutRequest.Inc()

	token, err := c.Cookie(keyRefreshToken)
	if err != nil {
		errs.SetBothErrors(c, models.HttpUnauthorized, err)
		return
	}

	if err = ctl.service.User().Logout(c, token); err != nil {
		errs.SetGinError(c, err)
		return
	}

	ctl.setRefreshToken(c, "")
	c.Status(http.StatusOK)
}

func (ctl *AuthController) VerifySession(c *gin.Context) {
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
		code := httpError.StatusCode
		log.Error(fmt.Sprintf("ValidateAccessToken: %s", err))

		if code == http.StatusInternalServerError {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, userSession)
}

func (ctl *AuthController) getAccessToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}

func (ctl *AuthController) setRefreshToken(ctx *gin.Context, token string) {
	//u, _ := url.Parse(ctl.config.ClientCallBackURI)

	ctx.SetCookie(keyRefreshToken,
		token,
		int(consts.RefreshTokenTTL),
		"/",
		"localhost",
		false,
		true)
}
