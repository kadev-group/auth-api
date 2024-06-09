package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type MobileController struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService
}

func InitMobileController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *MobileController {
	return &MobileController{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

func (ctl *MobileController) SignIn(c *gin.Context) {
	ctl.metrics.SignInRequests.Inc()

	var request models.MobileSignInReq
	if err := c.ShouldBindJSON(&request); err != nil {
		errs.SetBothErrors(c, models.HttpBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		errs.SetGinError(c, err)
		return
	}

	response, err := ctl.service.User().Authenticate(c, request.ToAuthReq())
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
