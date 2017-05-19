package helix

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	Conf      string
	Port      string
	Debug     bool
	Logfile   string
	Root      string
	StaticDir string
	Hostname  string
	Cert      string
	Key       string
}

func NewConfig() *Config {
	return &Config{
		Port:    "8443",
		Root:    GetCurrentRoot(),
		Debug:   false,
		Logfile: "",
		Cert:    "test_cert.pem",
		Key:     "test_key.pem",
	}
}

// LoadJSONFile loads server configuration
func (c *Config) LoadJSONFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &c)
}

func GetCurrentRoot() string {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	return root
}
