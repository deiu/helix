package helix

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OPTIONS(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", testServer.URL+"/foo", nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
