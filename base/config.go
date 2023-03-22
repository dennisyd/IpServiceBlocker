package base

type Configuration struct {
	IPBLOCKERSERVICE_ENVIRONMENT string

	// Defines which API type to use [HTTP | gRPC]
	IPBLOCKER_API_TYPE string

	IPBLOCKERSERVICE_SERVER_IP_ADDRESS string
	IPBLOCKERSERVICE_SERVER_PORT       string

	//
	IPBLOCKERSERVICE_USE_GCS bool

	// Defines location of the GeoLite2 binary database on disk
	IPBLOCKERSERVICE_GEOLITE2_PATH string
}

func GetConfig() Configuration {
	config := Configuration{}

	config.IPBLOCKERSERVICE_ENVIRONMENT = GetEnvAsStringOrDefault("IPBLOCKERSERVICE_ENVIRONMENT", DEV_ENVIRONMENT)
	config.IPBLOCKER_API_TYPE = GetEnvAsStringOrDefault("IPBLOCKER_API_TYPE", "HTTP")
	config.IPBLOCKERSERVICE_SERVER_IP_ADDRESS = GetEnvAsStringOrDefault("IPBLOCKERSERVICE_SERVER_IP_ADDRESS", "0.0.0.0")
	config.IPBLOCKERSERVICE_SERVER_PORT = GetEnvAsStringOrDefault("IPBLOCKERSERVICE_SERVER_PORT", "8080")
	config.IPBLOCKERSERVICE_GEOLITE2_PATH = GetEnvAsStringOrDefault("IPBLOCKERSERVICE_GEOLITE2_PATH", "GeoLite2-Country.mmdb")

	config.IPBLOCKERSERVICE_USE_GCS = GetEnvAsBoolOrDefault("IPBLOCKERSERVICE_USE_GCS", false)

	return config
}
