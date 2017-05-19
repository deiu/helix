package helix

import (
	"crypto/tls"
	"errors"
	"os"
	"path"

	"github.com/gocraft/web"
	"github.com/rs/zerolog"
)

var (
	HELIX_VERSION = "0.1"
	methodsAll    = []string{
		"OPTIONS", "HEAD", "GET", "POST", "PUT", "PATCH", "DELETE",
	}
	logger zerolog.Logger
)

type (
	GraphStore struct {
		Subject   map[string]string
		Predicate map[string]string
		Object    map[string]string
		Graph     map[string]string
		Version   string
		Current   bool
	}

	Context struct {
		Body   string
		Config *Config
		Store  *GraphStore
	}
)

func NewServer(config *Config) *web.Router {
	ctx := Context{}
	ctx.Config = config
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	if config.Debug == false {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	currentRoot, _ := os.Getwd()
	config.StaticDir = path.Join(currentRoot, config.StaticDir)

	// Middleware(web.StaticMiddleware(StaticRoot, web.StaticOption{Prefix: "/assets/"})).
	router := web.New(ctx). // Create your router
				Middleware((*Context).RequestLogger).
				Middleware(web.ShowErrorsMiddleware). //
				OptionsHandler((*Context).OptionsHandler).
				Get("/", (&ctx).RootHandler).
				Get("/:*", (&ctx).GetHandler). // Add a route
				Post("/:*", (&ctx).PostHandler).
				Put("/:*", (&ctx).PutHandler).
				Delete("/:*", (&ctx).DeleteHandler).
				Patch("/:*", (&ctx).PatchHandler).
				NotFound((*Context).NotFound)
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
