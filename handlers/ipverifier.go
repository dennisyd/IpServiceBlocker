package handlers

import (
	"context"
	"fmt"
	"net"
	"strings"

	base "IPBlockerService/base"
	data "IPBlockerService/data"
)

func ValidateIPAddress(ctx context.Context, ipAddress net.IP, validCountries []string) (bool, error) {
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

func IsDbAvailable(ctx context.Context) bool {
	databaseAccessable, _ := data.IsDbAvailable(ctx) // Error logged in data layer
	return databaseAccessable
}
