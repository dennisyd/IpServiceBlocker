package unit_tests

import (
	"log"
	"os"
)

func Setup() {
	log.Print(os.Getwd())
	os.Setenv("IPBLOCKERSERVICE_GEOLITE2_PATH", "../../data/GeoLite2-Country.mmdb")
}

func Teardown() {
	os.Unsetenv("IPBLOCKERSERVICE_GEOLITE2_PATH")
}
