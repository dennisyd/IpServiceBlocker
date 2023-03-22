package apishttp

import (
	"IPBlockerService/base"
	"IPBlockerService/errors"
	"IPBlockerService/handlers"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	ServiceOnline      bool `json:"serviceOnline"`
	DatabaseAccessable bool `json:"databaseAccessable"`
}

type VerifyIPRequest struct {
	IPAddress      string   `json:"ipAddress"`
	ValidCountries []string `json:"validCountries"`
}

type VerifyIPResponse struct {
	Authorized bool `json:"authorized"`
}

type HTTPAPI struct {
	Router *gin.Engine
}

func (api *HTTPAPI) Run(address string) {
	api.Router.Run(address)
}

func (api *HTTPAPI) healthCheck(ctx *gin.Context) {
	logWhenFinished, functionName := base.HandleNewMessage(ctx)
	defer logWhenFinished(ctx, functionName)

	res := HealthCheckResponse{
		ServiceOnline:      true,
		DatabaseAccessable: handlers.VerifyDbAvailable(ctx),
	}
	base.LogResponse(ctx, res)

	ctx.JSON(http.StatusOK, res)
}

func (api *HTTPAPI) verifyIP(ctx *gin.Context) {
	var newRequest VerifyIPRequest
	logWhenFinished, functionName := base.HandleNewMessage(ctx)
	defer logWhenFinished(ctx, functionName)

	if err := ctx.BindJSON(&newRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, "Improperly formed request.")
		return
	}

	ip := net.ParseIP(newRequest.IPAddress)
	if ip == nil {
		err := errors.InvalidIPAddressError{}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	isValid, err := handlers.ValidateIPAddress(ctx, ip, newRequest.ValidCountries)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	res := VerifyIPResponse{
		Authorized: isValid,
	}
	base.LogResponse(ctx, res)

	ctx.JSON(http.StatusOK, res)
}

func SetupHTTPAPI() *HTTPAPI {
	router := gin.Default()

	api := HTTPAPI{
		Router: router,
	}

	api.Router.GET("/healthcheck", api.healthCheck)
	api.Router.POST("/verifyip", api.verifyIP)

	config := base.GetConfig()
	if config.IPBLOCKERSERVICE_ENVIRONMENT == base.PROD_ENVIRONMENT {
		gin.SetMode(gin.ReleaseMode)
	}

	return &api
}
