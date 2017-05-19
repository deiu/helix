package helix

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Root(t *testing.T) {
	request, err := http.NewRequest("GET", testServer.URL, nil)
	assert.NoError(t, err)
	response, err := testClient.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	request, err = http.NewRequest("GET", testServer.URL+"/", nil)
	assert.NoError(t, err)
	response, err = testClient.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
}
