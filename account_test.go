package helix

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetAccountHandler(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+"/account/", nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func Test_GetAccountHandlerBadAuthz(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+"/account/", nil)
	assert.NoError(t, err)
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	req, err = http.NewRequest("GET", testServer.URL+"/account/", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", res.Header.Get("WWW-Authenticate"))
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func Test_AccountIntegration(t *testing.T) {
	config := NewConfig()
	config.StaticDir = testDir

	boltpath, err := newTempFile(config.StaticDir, "tmpbolt")
	assert.NoError(t, err)
	config.BoltPath = boltpath

	err = config.StartBolt()
	assert.NoError(t, err)

	ts, err := newTestServer(config)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", ts.URL+"/account/new", nil)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	form := url.Values{}
	form.Add("username", testUser)

	req, err = http.NewRequest("POST", ts.URL+"/account/new", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	form.Add("password", testPass)
	form.Add("email", testEmail)

	req, err = http.NewRequest("POST", ts.URL+"/account/new", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/login", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, testUser, res.Header.Get("User"))

	token := res.Header.Get("Token")
	assert.NotEmpty(t, token)

	req, err = http.NewRequest("GET", ts.URL+"/account/", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/logout", nil)
	assert.NoError(t, err)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/logout", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer")
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/logout", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/delete", nil)
	assert.NoError(t, err)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/delete", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL+"/account/delete", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	boltCleanup(config)
}

func Test_LoginBad(t *testing.T) {
	form := url.Values{}
	form.Add("username", testUser)

	req, err := http.NewRequest("POST", testServer.URL+"/account/login", nil)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Empty(t, res.Header.Get("User"))

	req, err = http.NewRequest("POST", testServer.URL+"/account/login", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err = testClient.Do(req)
	assert.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Empty(t, res.Header.Get("User"))
}

func Test_AddUser(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.addUser("", "", "")
	assert.Error(t, err)

	err = ctx.addUser(testUser, "", "")
	assert.Error(t, err)

	err = ctx.addUser(testUser, testPass, "")
	assert.Error(t, err)

	err = ctx.addUser("", testPass, "")
	assert.Error(t, err)

	err = ctx.addUser("", "", testEmail)
	assert.Error(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.Error(t, err)

	err = ctx.Config.StartBolt()
	assert.NoError(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	user, err := ctx.getUser(testUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user.Username)
	assert.Equal(t, testEmail, user.Email)

	boltCleanup(ctx.Config)
}

func Test_DeleteUser(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.Config.StartBolt()
	assert.NoError(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	err = ctx.deleteUser(testUser)
	assert.NoError(t, err)

	boltCleanup(ctx.Config)
}

func Test_SaveUserFail(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.Config.StartBolt()
	assert.NoError(t, err)

	user := NewUser()

	err = ctx.saveUser(user)
	assert.Error(t, err)

	boltCleanup(ctx.Config)
}

func Test_GetUser(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.Config.StartBolt()
	assert.NoError(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	_, err = ctx.getUser("foo")
	assert.Error(t, err)

	user, err := ctx.getUser(testUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user.Username)
	assert.Equal(t, testEmail, user.Email)

	boltCleanup(ctx.Config)
}

func boltCleanup(cfg *Config) {
	cfg.BoltDB.Close()
	err := os.Remove(cfg.BoltPath)
	if err != nil {
		panic("Failed to remove " + cfg.BoltPath)
	}
}
