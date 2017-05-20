package helix

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTLSConfig_NoCertKey(t *testing.T) {
	_, err := NewTLSConfig("", "test_key.pem")
	assert.Error(t, err)

	_, err = NewTLSConfig("test_cert.pem", "")
	assert.Error(t, err)
}

func Test_HTTP11(t *testing.T) {
	// Create a temporary http/1.1 client
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"http/1.1"},
			},
		},
	}
	req, err := http.NewRequest("GET", testServer.URL, nil)
	assert.NoError(t, err)

	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, res.ProtoAtLeast(1, 1))
}

func Test_HTTP2(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL, nil)
	assert.NoError(t, err)

	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, res.ProtoAtLeast(2, 0))
}
