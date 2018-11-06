package main

import (
	"encoding/json"
	"os"
)

type microConfig struct {
	Address      string
	Port         string
	ServeDir     string
	ServeIndex   string
	Errors       map[string]string
	Headers      map[string]string
	Methods      map[string]string
	ContentTypes contentTypes
	Proxy        proxy
	TLS          bool
	TLSCert      string
	TLSKey       string
}

type proxy struct {
	Enabled bool
	Rules   map[string]string
}

type contentTypes struct {
	ResponseTypes map[string]string
	RequestTypes  []string
}

func loadConfigFromFile(p string, c *microConfig) {

	// check if config exists
	if _, err := os.Stat(p); err != nil {
		logAction(logERROR, err)
	}

	// load config
	file, err := os.Open(p)
	logAction(logERROR, err)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	logAction(logERROR, err)
}
