package helix

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PATCH(t *testing.T) {
	req, err := http.NewRequest("PATCH", testServer.URL+"/foo", nil)
	assert.NoError(t, err)
	resp, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
