package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"auth-api/internal/pkg/tools"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	metrics *metrics.APIMetrics
	log     *zap.Logger
	config  *models.Config
	service interfaces.IService
}

func InitUserController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *UserController {
	return &UserController{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

func (ctl *UserController) GetByUserIDCode(c *gin.Context) {
	ctl.metrics.VerifyEmailRequests.Inc()

	userIDCode := c.Param("user_idcode")
	if !tools.IsUUID(userIDCode) {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	response, err := ctl.service.User().GetByUserIDCode(c, userIDCode)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
