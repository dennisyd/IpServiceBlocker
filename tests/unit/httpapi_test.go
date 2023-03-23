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

func TestHealthCheckRoute_http(t *testing.T) {
	Setup()
	defer Teardown()
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	api.ServeHTTP(w, req)

	data := apishttp.HealthCheckResponse{}

	assert.Equal(t, 200, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, data.ServiceOnline, true)
	assert.Equal(t, data.DatabaseAccessable, true)
}

func TestHealthCheckRoute_http_db_inaccessable(t *testing.T) {
	// Don't perform setup/teardown so the db file cannot be found
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	api.ServeHTTP(w, req)

	data := apishttp.HealthCheckResponse{}

	assert.Equal(t, 200, w.Code)
	err := json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, data.ServiceOnline, true)
	assert.Equal(t, data.DatabaseAccessable, false)
}

func TestAuthroizeIPRoute_http(t *testing.T) {
	Setup()
	defer Teardown()
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()

	requestData := apishttp.AuthroizeIPRequest{
		IPAddress: "162.226.203.50",
		ValidCountries: []string{
			"United States",
			"United Kingdom",
			"Canada",
		},
	}

	jsonRequestData, err := json.Marshal(requestData)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/authorizeip", bytes.NewReader(jsonRequestData))
	api.ServeHTTP(w, req)

	data := apishttp.AuthroizeIPResponse{}

	assert.Equal(t, 200, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, true, data.Authorized)
}

func TestAuthroizeIPRoute_http_invalid_ip(t *testing.T) {
	Setup()
	defer Teardown()
	api := apishttp.SetupHTTPAPI()
	w := httptest.NewRecorder()

	requestData := apishttp.AuthroizeIPRequest{
		IPAddress: "103.136.43.2",
		ValidCountries: []string{
			"United States",
			"United Kingdom",
			"Canada",
		},
	}

	jsonRequestData, err := json.Marshal(requestData)
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/authorizeip", bytes.NewReader(jsonRequestData))
	api.ServeHTTP(w, req)

	data := apishttp.AuthroizeIPResponse{}

	assert.Equal(t, 200, w.Code)
	err = json.Unmarshal([]byte(w.Body.String()), &data)
	assert.Nil(t, err)
	assert.Equal(t, false, data.Authorized)
}
