package unit_tests

import (
	handlers "IPBlockerService/handlers"
	"net"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateIPAddress_no_valid_countries(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := gin.Context{}
	res, err := handlers.ValidateIPAddress(&ctx, net.ParseIP("81.2.69.142"), []string{})
	assert.Nil(t, err)
	assert.False(t, res)

	res, err = handlers.ValidateIPAddress(&ctx, net.ParseIP("103.159.28.2"), []string{})
	assert.Nil(t, err)
	assert.False(t, res)

	res, err = handlers.ValidateIPAddress(&ctx, net.ParseIP("1.1.1.1"), []string{})
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestValidateIPAddress_valid_countries(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := gin.Context{}
	res, err := handlers.ValidateIPAddress(&ctx, net.ParseIP("81.2.69.142"), []string{
		"United Kingdom",
	})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = handlers.ValidateIPAddress(&ctx, net.ParseIP("103.159.28.2"), []string{
		"Palau",
	})
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestValidateIPAddress_invalid_countries(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := gin.Context{}
	res, err := handlers.ValidateIPAddress(&ctx, net.ParseIP("103.136.43.2"), []string{
		"United States",
		"United Kingdom",
		"Palau",
	})
	assert.Nil(t, err)
	assert.False(t, res)

	res, err = handlers.ValidateIPAddress(&ctx, net.ParseIP("102.129.166.45"), []string{
		"United States",
		"United Kingdom",
		"Palau",
	})
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestValidateIPAddress_invalid_ips(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := gin.Context{}
	res, err := handlers.ValidateIPAddress(&ctx, net.ParseIP("adwawd"), []string{
		"United States",
		"United Kingdom",
		"Palau",
	})
	assert.NotNil(t, err)
	assert.False(t, res)

	res, err = handlers.ValidateIPAddress(&ctx, net.ParseIP("11.1.2"), []string{
		"United States",
		"United Kingdom",
		"Palau",
	})
	assert.NotNil(t, err)
	assert.False(t, res)
}
