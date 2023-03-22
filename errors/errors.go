// Holds all possible error definitions for the IPBlockerService

package errors

type InvalidIPAddressError struct{}

func (e *InvalidIPAddressError) Error() string {
	return "Received an invalid IP address"
}

type GeoLite2DBError struct{}

func (e *GeoLite2DBError) Error() string {
	return "Error while reading the GeoLite2 DB"
}
