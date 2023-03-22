package unit_tests

import (
	apishttp "IPBlockerService/apis/http"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoute(t *testing.T) {
	Setup()
	defer Teardown()
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	api.Router.ServeHTTP(w, req)

	data := apishttp.HealthCheckResponse{}

	assert.Equal(t, 200, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, data.ServiceOnline, true)
	assert.Equal(t, data.DatabaseAccessable, true)
}

func TestHealthCheckRoute_db_inaccessable(t *testing.T) {
	// Don't perform setup/teardown so the db file cannot be found
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	api.Router.ServeHTTP(w, req)

	data := apishttp.HealthCheckResponse{}

	assert.Equal(t, 200, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, data.ServiceOnline, true)
	assert.Equal(t, data.DatabaseAccessable, false)
}

func TestVerifyIPRoute_invalid_ip(t *testing.T) {
	Setup()
	defer Teardown()
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()

	requestData := apishttp.VerifyIPRequest{
		// IPAddress: "162.226.203.50",
		// IPAddress: "100.42.20.234",
		IPAddress: "2600:1700:2890:d700:848:aa7c:99d2:8914",
		ValidCountries: []string{
			"United States",
			"United Kingdom",
			"Canada",
		},
	}

	jsonRequestData, err := json.Marshal(requestData)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/verifyip", bytes.NewReader(jsonRequestData))
	api.Router.ServeHTTP(w, req)

	data := apishttp.VerifyIPResponse{}

	assert.Equal(t, 200, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, true, data.Authorized)
}
