package helix

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type HelixConfig struct {
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

func NewHelixConfig() *HelixConfig {
	return &HelixConfig{
		Port:    "8443",
		Root:    GetCurrentRoot(),
		Debug:   false,
		Logfile: "",
		Cert:    "test_cert.pem",
		Key:     "test_key.pem",
	}
}

// LoadJSONFile loads server configuration
func (c *HelixConfig) LoadJSONFile(filename string) error {
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

func (c *HelixConfig) GetLogger() io.Writer {
	if len(c.Logfile) == 0 {
		return os.Stderr
	}
	w, err := os.OpenFile(c.Logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return w
}
