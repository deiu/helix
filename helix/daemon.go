package main

import (
	"net/http"
	"os"

	"github.com/deiu/helix"
)

var (
	port   = os.Getenv("HELIX_PORT")
	host   = os.Getenv("HELIX_HOST")
	root   = os.Getenv("HELIX_ROOT")
	static = os.Getenv("HELIX_STATIC_DIR")
	debug  = os.Getenv("HELIX_DEBUG")
	log    = os.Getenv("HELIX_LOGGING")
	cert   = os.Getenv("HELIX_CERT")
	key    = os.Getenv("HELIX_KEY")
	hsts   = os.Getenv("HELIX_HSTS")
	bolt   = os.Getenv("HELIX_BOLT_PATH")
)

func main() {
	println("Starting server...")

	config := helix.NewConfig()
	config.Port = port
	config.Hostname = host
	config.Root = root
	config.StaticDir = static
	config.Cert = cert
	config.Key = key
	config.BoltPath = bolt
	if len(debug) > 0 {
		config.Debug = true
	}
	if len(log) > 0 {
		config.Logging = true
	}
	if len(hsts) > 0 {
		config.HSTS = true
	}

	println("Listening on " + config.Hostname + ":" + config.Port)

	ctx := helix.NewContext()

	if len(config.BoltPath) > 0 {
		// Start Bolt
		err := ctx.StartBolt()
		if err != nil {
			println(err.Error())
			return
		}
		defer ctx.BoltDB.Close()
	}

	// prepare new handler
	handler := helix.NewServer(ctx)
	// prepare server
	s := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}
	// set TLS config
	tlsCfg, err := helix.NewTLSConfig(config.Cert, config.Key)
	if err != nil {
		println(err.Error())
		return
	}
	s.TLSConfig = tlsCfg
	// start server
	s.ListenAndServeTLS(config.Cert, config.Key)
}
