package main

import (
	"net/http"
	"os"

	helix "github.com/deiu/gold"
)

var (
	port   = os.Getenv("HELIX_PORT")
	host   = os.Getenv("HELIX_HOST")
	root   = os.Getenv("HELIX_ROOT")
	static = os.Getenv("HELIX_STATIC_DIR")
	debug  = os.Getenv("HELIX_DEBUG")
	cert   = os.Getenv("HELIX_CERT")
	key    = os.Getenv("HELIX_KEY")
)

func main() {
	println("Starting server...")
	config := helix.NewHelixConfig()
	config.Port = port
	config.Hostname = host
	config.Root = root
	config.StaticDir = static
	config.Cert = cert
	config.Key = key
	if len(debug) > 0 {
		config.Debug = true
	}

	router := helix.NewServer(config)
	http.ListenAndServe(":"+config.Port, router) // Start the server!
}
