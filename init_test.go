package helix

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/http2"
)

var (
	testCtx    *Context
	testServer *httptest.Server
	testClient *http.Client

	testUser  = "alice"
	testPass  = "testpass"
	testEmail = "foo@bar.baz"
	testDir   = "./test"
)

func init() {
	// uncomment for extra logging
	config := NewConfig()
	config.Debug = true
	config.HSTS = true
	config.StaticDir = testDir

	testCtx = NewContext()
	testCtx.Config = config

	var err error
	testServer, err = newTestServer(testCtx)
	if err != nil {
		println(err.Error())
		return
	}
	// testClient
	testClient = newTestClient()
}

func newTestServer(ctx *Context) (*httptest.Server, error) {
	var ts *httptest.Server
	// testServer
	handler := NewServer(ctx)
	ts = httptest.NewUnstartedServer(handler)

	// prepare TLS config
	tlsCfg, err := NewTLSConfig(ctx.Config.Cert, ctx.Config.Key)
	if err != nil {
		return ts, err
	}

	ts.TLS = tlsCfg
	ts.StartTLS()

	ts.URL = strings.Replace(ts.URL, "127.0.0.1", "localhost", 1)

	return ts, nil
}

func newTestClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"h2"},
			},
		},
	}
}

func newTempFile(dir, name string) (string, error) {
	tmpfile, err := ioutil.TempFile(dir, name)
	if err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}
