package helix

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"path"

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
		Config      *Config
		Store       map[string]*Graph
		User        string
		AccessToken string
	}
)

func NewContext() *Context {
	return &Context{
		Config: NewConfig(),
		Store:  make(map[string]*Graph),
		User:   "",
	}
}

func NewServer(cfg *Config) *web.Router {
	ctx := NewContext()
	ctx.Config = cfg
	if !ctx.Config.Logging {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	currentRoot, _ := os.Getwd()
	ctx.Config.StaticDir = path.Join(currentRoot, ctx.Config.StaticDir)

	// Create router and add middleware
	router := web.New(*ctx).
		// Middleware(web.LoggerMiddleware). // turn off once done with tweaking
		Middleware(ctx.CORS).
		Middleware(ctx.Authentication).
		Middleware(ctx.RequestLogger).
		Middleware(web.StaticMiddleware(ctx.Config.StaticDir, web.StaticOption{Prefix: ctx.Config.StaticPath})).
		OptionsHandler(ctx.OptionsHandler)

	// Account routes
	router.Get("/", ctx.RootHandler).
		Post("/account/new", ctx.NewAccountHandler).
		Post("/account/logout", ctx.LogoutHandler).
		Post("/account/login", ctx.LoginHandler).
		Post("/account/delete", ctx.DeleteAccountHandler).
		Get("/account/", ctx.GetAccountHandler)

	// API routes
	router.Get("/:*", ctx.GetHandler).
		Post("/:*", ctx.PostHandler).
		Put("/:*", ctx.PutHandler).
		Delete("/:*", ctx.DeleteHandler).
		Patch("/:*", ctx.PatchHandler)

	if ctx.Config.Debug {
		router.Middleware(web.ShowErrorsMiddleware)
	}

	return router
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
