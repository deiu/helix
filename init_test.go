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
	testConfig *Config
	testServer *httptest.Server
	testClient *http.Client

	testUser  = "alice"
	testPass  = "testpass"
	testEmail = "foo@bar.baz"
	testDir   = "./test"
)

func init() {
	// uncomment for extra logging
	testConfig = NewConfig()
	testConfig.Debug = true
	testConfig.HSTS = true
	testConfig.StaticDir = testDir

	var err error
	testServer, err = newTestServer(testConfig)
	if err != nil {
		println(err.Error())
		return
	}
	// testClient
	testClient = newTestClient()
}

func newTestServer(cfg *Config) (*httptest.Server, error) {
	var ts *httptest.Server
	// testServer
	handler := NewServer(cfg)
	ts = httptest.NewUnstartedServer(handler)

	// prepare TLS config
	tlsCfg, err := NewTLSConfig(cfg.Cert, cfg.Key)
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
