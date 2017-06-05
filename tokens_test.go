package helix

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gocraft/web"
	"github.com/stretchr/testify/assert"
)

type testWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func Test_GetAuthzUserFromToken(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()
	ctx.Config.StaticDir = testDir

	boltpath, err := newTempFile(ctx.Config.StaticDir, "tmpbolt")
	assert.NoError(t, err)
	ctx.Config.BoltPath = boltpath

	err = ctx.Config.StartBolt()
	assert.NoError(t, err)

	tokenType := "Authorization"
	host := "localhost"
	origin := "example.org"
	values := map[string]string{
		"webid": testUser,
	}

	token, err := ctx.newPersistedToken(tokenType, host, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ctx.getAuthzUserFromToken("", host, origin)
	assert.Error(t, err)

	_, err = ctx.getAuthzUserFromToken(token, "", origin)
	assert.Error(t, err)

	_, err = ctx.getAuthzUserFromToken(token, host, "")
	assert.Error(t, err)

	webid, err := ctx.getAuthzUserFromToken(token, host, "foo")
	assert.Error(t, err)
	assert.Empty(t, webid)

	webid, err = ctx.getAuthzUserFromToken(token, host, origin)
	assert.Error(t, err)
	assert.Empty(t, webid)

	values["origin"] = "example.org"
	token, err = ctx.newPersistedToken(tokenType, host, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ctx.getAuthzUserFromToken(token, host, "foo")
	assert.Error(t, err)

	webid, err = ctx.getAuthzUserFromToken(token, host, origin)
	assert.NoError(t, err)
	assert.Equal(t, testUser, webid)

	values["valid"] = fmt.Sprintf("%d", time.Now().Add(time.Duration(1)*time.Microsecond).UnixNano())
	token, err = ctx.newPersistedToken(tokenType, host, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	time.Sleep(time.Millisecond * 1)

	webid, err = ctx.getAuthzUserFromToken(token, host, origin)
	assert.Error(t, err)
	assert.Empty(t, webid)

	boltCleanup(ctx.Config)
}

func Test_PersistedTokens(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()
	ctx.Config.StaticDir = testDir

	boltpath, err := newTempFile(ctx.Config.StaticDir, "tmpbolt")
	assert.NoError(t, err)
	ctx.Config.BoltPath = boltpath

	err = ctx.Config.StartBolt()
	assert.NoError(t, err)

	tokenType := "Authorization"
	host := "localhost"
	origin := "example.org"
	values := map[string]string{
		"webid":  testUser,
		"origin": origin,
	}

	_, err = ctx.newPersistedToken("", host, values)
	assert.Error(t, err)

	_, err = ctx.newPersistedToken(tokenType, "", values)
	assert.Error(t, err)

	token, err := ctx.newPersistedToken(tokenType, host, values)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ctx.getPersistedToken("", host, token)
	assert.Error(t, err)

	_, err = ctx.getPersistedToken(tokenType, "", token)
	assert.Error(t, err)

	_, err = ctx.getPersistedToken(tokenType, host, "")
	assert.Error(t, err)

	_, err = ctx.getPersistedToken("bar", host, token)
	assert.Error(t, err)

	_, err = ctx.getPersistedToken(tokenType, "foo", token)
	assert.Error(t, err)

	vals, err := ctx.getPersistedToken(tokenType, host, token)
	assert.NoError(t, err)
	assert.Equal(t, values["webid"], vals["webid"])
	assert.Equal(t, values["origin"], vals["origin"])

	_, err = ctx.getTokenByOrigin("", host, origin)
	assert.Error(t, err)

	_, err = ctx.getTokenByOrigin(tokenType, "", origin)
	assert.Error(t, err)

	_, err = ctx.getTokenByOrigin(tokenType, host, "")
	assert.Error(t, err)

	_, err = ctx.getTokenByOrigin("foo", host, origin)
	assert.Error(t, err)

	_, err = ctx.getTokenByOrigin(tokenType, "bar", origin)
	assert.Error(t, err)

	tkn, err := ctx.getTokenByOrigin(tokenType, host, "test.com")
	assert.NoError(t, err)
	assert.Empty(t, tkn)

	tkn, err = ctx.getTokenByOrigin(tokenType, host, origin)
	assert.NoError(t, err)
	assert.Equal(t, token, tkn)

	_, err = ctx.getTokensByType("", host)
	assert.Error(t, err)

	_, err = ctx.getTokensByType("foo", host)
	assert.Error(t, err)

	_, err = ctx.getTokensByType(tokenType, "baz")
	assert.Error(t, err)

	_, err = ctx.getTokensByType(tokenType, "")
	assert.Error(t, err)

	tokens, err := ctx.getTokensByType(tokenType, host)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(tokens))
	assert.Equal(t, tokens[token]["webid"], values["webid"])
	assert.Equal(t, tokens[token]["origin"], values["origin"])

	err = ctx.deletePersistedToken("", host, token)
	assert.Error(t, err)

	err = ctx.deletePersistedToken("foo", host, token)
	assert.Error(t, err)

	err = ctx.deletePersistedToken(tokenType, "", token)
	assert.Error(t, err)

	err = ctx.deletePersistedToken(tokenType, "foo", token)
	assert.Error(t, err)

	err = ctx.deletePersistedToken(tokenType, host, "")
	assert.Error(t, err)

	err = ctx.deletePersistedToken(tokenType, host, token)
	assert.NoError(t, err)

	boltCleanup(ctx.Config)
}

func Test_NewAuthzToken(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()
	ctx.Config.StaticDir = testDir

	boltpath, err := newTempFile(ctx.Config.StaticDir, "tmpbolt")
	assert.NoError(t, err)
	ctx.Config.BoltPath = boltpath

	err = ctx.Config.StartBolt()
	assert.NoError(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", testServer.URL, nil)
	rec := httptest.NewRecorder()
	w := web.ResponseWriter(&testWriter{ResponseWriter: rec})

	ctx.newAuthzToken(w, req, testUser)

	assert.Equal(t, w.Header().Get("Token"), rec.Header().Get("Token"))

	ctx.Config.BoltDB.Close()

	rec = httptest.NewRecorder()
	w = web.ResponseWriter(&testWriter{ResponseWriter: rec})

	ctx.newAuthzToken(w, req, testUser)
	assert.Empty(t, rec.Header().Get("Token"))

	boltCleanup(ctx.Config)
}

func Test_TokenDateIsValid(t *testing.T) {
	err := tokenDateIsValid("")
	assert.Error(t, err)

	valid := fmt.Sprintf("%d", time.Now().Add(time.Duration(1)*time.Microsecond).UnixNano())
	time.Sleep(time.Millisecond * 1)

	err = tokenDateIsValid(valid)
	assert.Error(t, err)

	valid = fmt.Sprintf("%d", time.Now().Add(time.Duration(1)*time.Millisecond).UnixNano())
	time.Sleep(time.Microsecond * 1)

	err = tokenDateIsValid(valid)
	assert.NoError(t, err)
}

// Don't need this yet because we get it for free:
func (w *testWriter) Write(data []byte) (n int, err error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(data)
	w.size += size
	return size, err
}

func (w *testWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *testWriter) StatusCode() int {
	return w.statusCode
}

func (w *testWriter) Written() bool {
	return w.statusCode != 0
}

func (w *testWriter) Size() int {
	return w.size
}

func (w *testWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (w *testWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *testWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
