package helix

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/boltdb/bolt"
	"github.com/gocraft/web"
	"github.com/rs/zerolog"
)

const HelixVersion = "0.1"

var (
	methodsAll = []string{
		"OPTIONS", "HEAD", "GET", "POST", "PUT", "PATCH", "DELETE",
	}
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
)

type (
	Context struct {
		Config *Config
		Store  map[string]*Graph
		BoltDB *bolt.DB
	}
)

func NewContext() *Context {
	return &Context{
		Config: NewConfig(),
		Store:  make(map[string]*Graph),
		BoltDB: &bolt.DB{},
	}
}

func NewServer(config *Config) (*web.Router, error) {
	ctx := NewContext()
	ctx.Config = config

	if !config.Logging {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	if len(config.BoltPath) > 0 {
		// Start Bolt
		err := ctx.StartBolt()
		if err != nil {
			return &web.Router{}, err
		}
		defer ctx.BoltDB.Close()
	}

	currentRoot, _ := os.Getwd()
	config.StaticDir = path.Join(currentRoot, config.StaticDir)

	// Create router
	router := web.New(*ctx).
		Middleware((ctx).RequestLogger).
		Middleware((ctx).CORSMiddleware).
		Middleware(web.StaticMiddleware(config.StaticDir, web.StaticOption{Prefix: config.StaticPath})).
		OptionsHandler((ctx).OptionsHandler).
		Get("/:*", (ctx).GetHandler). // Match anything
		Post("/:*", (ctx).PostHandler).
		Put("/:*", (ctx).PutHandler).
		Delete("/:*", (ctx).DeleteHandler).
		Patch("/:*", (ctx).PatchHandler).
		Get("/", (ctx).RootHandler) // reserved for a welcome/dashboard page

	if config.Debug {
		router.Middleware(web.ShowErrorsMiddleware)
	}

	return router, nil
}

func NewTLSConfig(cert, key string) (*tls.Config, error) {
	var err error
	cfg := &tls.Config{}

	if len(cert) == 0 || len(key) == 0 {
		return cfg, errors.New("Missing cert and key for TLS configuration")
	}

	cfg.MinVersion = tls.VersionTLS12
	cfg.NextProtos = []string{"h2"}
	// use strong crypto
	cfg.PreferServerCipherSuites = true
	cfg.CurvePreferences = []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256}
	cfg.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	}
	cfg.Certificates = make([]tls.Certificate, 1)
	cfg.Certificates[0], err = tls.LoadX509KeyPair(cert, key)

	return cfg, err
}

func absoluteURI(req *http.Request) string {
	scheme := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		scheme += "s"
	}
	reqHost := req.Host
	if len(req.Header.Get("X-Forward-Host")) > 0 {
		reqHost = req.Header.Get("X-Forward-Host")
	}
	host, port, err := net.SplitHostPort(reqHost)
	if err != nil {
		host = reqHost
	}
	if len(host) == 0 {
		host = "localhost"
	}
	if len(port) > 0 {
		port = ":" + port
	}
	if (scheme == "https" && port == ":443") || (scheme == "http" && port == ":80") {
		port = ""
	}
	return scheme + "://" + host + port + req.URL.Path
}

func (c *Context) StartBolt() error {
	var err error
	c.BoltDB, err = bolt.Open(c.Config.BoltPath, 0644, nil)
	if err != nil {
		return err
	}
	return nil
}
