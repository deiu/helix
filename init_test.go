package helix

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/http2"
)

var (
	testConfig *Config
	testServer *httptest.Server
	testClient *http.Client
)

func init() {
	// uncomment for extra logging
	testConfig := NewConfig()
	testConfig.Debug = true
	testConfig.HSTS = true
	testConfig.RedisURL = "localhost:1234"
	testConfig.StaticDir = "./static"

	// prepare TLS config
	tlsCfg, err := NewTLSConfig(testConfig.Cert, testConfig.Key)
	if err != nil {
		println(err.Error())
		return
	}

	// testServer
	e := NewServer(testConfig)
	testServer = httptest.NewUnstartedServer(e)
	testServer.TLS = tlsCfg
	testServer.StartTLS()

	testServer.URL = strings.Replace(testServer.URL, "127.0.0.1", "localhost", 1)
	// testClient
	testClient = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"h2"},
			},
		},
	}
}
