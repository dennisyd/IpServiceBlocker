package apishttp

import (
	"IPBlockerService/base"
	"IPBlockerService/errors"
	"IPBlockerService/handlers"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	ServiceOnline      bool `json:"serviceOnline"`
	DatabaseAccessable bool `json:"databaseAccessable"`
}

type AuthroizeIPRequest struct {
	IPAddress      string   `json:"ipAddress"`
	ValidCountries []string `json:"validCountries"`
}

type AuthroizeIPResponse struct {
	Authorized bool `json:"authorized"`
}

func addMsgId(ctx *gin.Context) {
	ctx.Set(base.CONTEXT_MESSAGE_ID, guuid.New())
}

func healthCheck(ctx *gin.Context) {
	addMsgId(ctx)
	logWhenFinished, functionName := base.HandleNewMessage(ctx)
	defer logWhenFinished(ctx, functionName)

	res := HealthCheckResponse{
		ServiceOnline:      true,
		DatabaseAccessable: handlers.IsDbAvailable(ctx),
	}
	base.LogResponse(ctx, res)

	ctx.JSON(http.StatusOK, res)
}

func authorizeIP(ctx *gin.Context) {
	addMsgId(ctx)
	var newRequest AuthroizeIPRequest

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

	res := AuthroizeIPResponse{
		Authorized: isValid,
	}
	base.LogResponse(ctx, res)

	ctx.JSON(http.StatusOK, res)
}

func SetupHTTPAPI() *gin.Engine {
	router := gin.Default()

	router.GET("/healthcheck", healthCheck)
	router.POST("/authorizeip", authorizeIP)

	config := base.GetConfig()
	if config.IPBLOCKERSERVICE_ENVIRONMENT == base.PROD_ENVIRONMENT {
		gin.SetMode(gin.ReleaseMode)
	}

	return router
}
