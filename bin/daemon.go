package main

import (
	"net/http"
	"os"

	helix "github.com/deiu/gold"
)

var (
	port    = os.Getenv("HELIX_PORT")
	host    = os.Getenv("HELIX_HOST")
	root    = os.Getenv("HELIX_ROOT")
	static  = os.Getenv("HELIX_STATIC_DIR")
	logfile = os.Getenv("HELIX_LOG")
	debug   = os.Getenv("HELIX_DEBUG")
	cert    = os.Getenv("HELIX_CERT")
	key     = os.Getenv("HELIX_KEY")
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
	config.Logfile = logfile
	if len(debug) > 0 {
		config.Debug = true
	}
	println("Listening on " + config.Hostname + ":" + config.Port)

	router := helix.NewServer(config)
	http.ListenAndServe(":"+config.Port, router) // Start the server!
}
