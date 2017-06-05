package helix

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

type Config struct {
	Conf       string
	Port       string
	Logging    bool
	Debug      bool
	SkipVerify bool
	Root       string
	StaticDir  string
	StaticPath string
	Hostname   string
	Cert       string
	Key        string
	TokenAge   int64
	HSTS       bool
	BoltPath   string
	BoltDB     *bolt.DB
	FilePath   string
	DataPath   string
	ACLPath    string
	MetaPath   string
}

func NewConfig() *Config {
	return &Config{
		Port:       "8443",
		Root:       GetCurrentRoot(),
		Logging:    false,
		Debug:      false,
		SkipVerify: false,
		Cert:       "test_cert.pem",
		Key:        "test_key.pem",
		TokenAge:   5,
		HSTS:       false,
		BoltPath:   filepath.Join(os.TempDir(), "bolt.db"),
		BoltDB:     &bolt.DB{},
		StaticPath: "/static/",
		FilePath:   "/files/",
		DataPath:   "/data/",
		ACLPath:    "/acl/",
		MetaPath:   "/meta/",
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

func (c *Config) StartBolt() error {
	var err error
	c.BoltDB, err = bolt.Open(c.BoltPath, 0664, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return nil
}
