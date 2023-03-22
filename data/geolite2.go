package data

import (
	"IPBlockerService/base"
	"IPBlockerService/errors"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

func GetOrignatingCountryFromIP(ctx *gin.Context, ipAddress net.IP) (string, error) {
	if ipAddress == nil {
		return "", &errors.InvalidIPAddressError{}
	}
	var db *geoip2.Reader
	var err error
	if db, err = getDb(ctx); err != nil {
		return "", err
	}

	defer db.Close()

	countryRecord, err := db.Country(ipAddress)
	if err != nil {
		base.PrintLogLine(ctx, err.Error())
		return "", &errors.GeoLite2DBError{}
	}

	if countryRecord == nil {
		base.PrintLogLine(ctx, fmt.Sprintf("GeoLite2 database returned a nil record for IP address \"%v\"", ipAddress.String()))
		return "", &errors.GeoLite2DBError{}
	}

	return countryRecord.Country.Names["en"], nil
}

func getDb(ctx *gin.Context) (*geoip2.Reader, error) {
	db, err := geoip2.Open(base.GetConfig().IPBLOCKERSERVICE_GEOLITE2_PATH)
	if err != nil {
		base.PrintLogLine(ctx, err.Error())
		return nil, &errors.GeoLite2DBError{}
	}
	return db, nil
}

func VerifyDbAvailable(ctx *gin.Context) (bool, error) {
	var db *geoip2.Reader
	var err error
	if db, err = getDb(ctx); err != nil {
		base.PrintLogLine(ctx, err.Error())
		return false, &errors.GeoLite2DBError{}
	}
	db.Close()

	return true, nil
}
