package handlers

import (
	"fmt"
	"net"
	"strings"

	base "IPBlockerService/base"
	data "IPBlockerService/data"

	"github.com/gin-gonic/gin"
)

func ValidateIPAddress(ctx *gin.Context, ipAddress net.IP, validCountries []string) (bool, error) {
	var originatingCountry string
	var err error

	if originatingCountry, err = data.GetOrignatingCountryFromIP(ctx, ipAddress); err != nil {
		return false, err
	}

	base.PrintLogLine(ctx, fmt.Sprintf("IP Address \"%v\" originated from country \"%v\"", ipAddress.String(), originatingCountry))

	for _, country := range validCountries {
		if strings.ToLower(country) == strings.ToLower(originatingCountry) {
			base.PrintLogLine(ctx, fmt.Sprintf("Originating country \"%v\" found in valid list", originatingCountry))
			return true, nil
		}
	}

	base.PrintLogLine(ctx, fmt.Sprintf("Originating country \"%v\" was not found in valid list", originatingCountry))
	return false, nil
}

func VerifyDbAvailable(ctx *gin.Context) bool {
	databaseAccessable, _ := data.VerifyDbAvailable(ctx) // Error logged in data layer
	return databaseAccessable
}
