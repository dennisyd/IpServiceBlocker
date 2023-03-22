package main

import (
	apishttp "IPBlockerService/apis/http"
	"IPBlockerService/base"
)

func main() {
	api := apishttp.SetupHTTPAPI()
	config := base.GetConfig()
	api.Run(config.IPBLOCKERSERVICE_SERVER_IP_ADDRESS + ":" + config.IPBLOCKERSERVICE_SERVER_PORT)
}
