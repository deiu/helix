package helix

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthenticationRequired(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+"/account/", nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func Test_ParseBearerAuthorizationHeader(t *testing.T) {
	token := "verylongnonce"
	h := "Bearer " + token
	tkn, err := ParseBearerAuthorizationHeader(h)
	assert.NoError(t, err)
	assert.Equal(t, token, tkn)

	tkn, err = ParseBearerAuthorizationHeader("")
	assert.Error(t, err)
	assert.Empty(t, tkn)

	h = "Foo bar"
	tkn, err = ParseBearerAuthorizationHeader(h)
	assert.Error(t, err)
	assert.Empty(t, tkn)
}

func Test_SavePassFail(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.savePass(testUser, "")
	assert.Error(t, err)

	err = ctx.Config.StartBolt()
	assert.NoError(t, err)

	err = ctx.savePass("foo", testPass)
	assert.Error(t, err)

	boltCleanup(ctx.Config)
}

func Test_VerifyPass(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.Config.StartBolt()
	assert.NoError(t, err)

	ok, err := ctx.verifyPass("", "")
	assert.Error(t, err)
	assert.False(t, ok)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	ok, err = ctx.verifyPass("bob", "foo")
	assert.Error(t, err)
	assert.False(t, ok)

	ok, err = ctx.verifyPass(testUser, "foo")
	assert.Error(t, err)
	assert.False(t, ok)

	ok, err = ctx.verifyPass(testUser, testPass)
	assert.NoError(t, err)
	assert.True(t, ok)

	boltCleanup(ctx.Config)
}
