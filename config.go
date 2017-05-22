package helix

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Conf      string
	Port      string
	Logging   bool
	Debug     bool
	Root      string
	StaticDir string
	Hostname  string
	Cert      string
	Key       string
	HSTS      bool
	FilePath  string
	DataPath  string
	ACLPath   string
	MetaPath  string
}

func NewConfig() *Config {
	return &Config{
		Port:     "8443",
		Root:     GetCurrentRoot(),
		Logging:  false,
		Debug:    false,
		Cert:     "test_cert.pem",
		Key:      "test_key.pem",
		HSTS:     false,
		FilePath: "/files/",
		DataPath: "/data/",
		ACLPath:  "/acl/",
		MetaPath: "/meta/",
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
	root, _ := os.Getwd()
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	return root
}
