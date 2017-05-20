package helix

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	exUri = "https://example.org"
)

func TestHSTS(t *testing.T) {
	req, err := http.NewRequest("HEAD", testServer.URL, nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, "max-age=63072000", res.Header.Get("Strict-Transport-Security"))
}

func Test_CORS_NoOrigin(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL, nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}

func Test_CORS_Origin(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL, nil)
	assert.NoError(t, err)

	req.Header.Set("Origin", exUri)
	res, err := testClient.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, exUri, res.Header.Get("Access-Control-Allow-Origin"))
}
func Test_CORS_AllowHeaders(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", testServer.URL, nil)
	assert.NoError(t, err)
	req.Header.Add("Access-Control-Request-Headers", "User, ETag")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Empty(t, string(body))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "User, ETag", res.Header.Get("Access-Control-Allow-Headers"))
}

func Test_CORS_NoReqMethod(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", testServer.URL, nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Empty(t, string(body))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, strings.Join(methodsAll, ", "), res.Header.Get("Access-Control-Allow-Methods"))
}

func Test_CORS_ReqMethod(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", testServer.URL, nil)
	assert.NoError(t, err)
	req.Header.Add("Access-Control-Request-Method", "PATCH")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Empty(t, string(body))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "PATCH", res.Header.Get("Access-Control-Allow-Methods"))
}
