package config

import "time"

type Config struct {
	Username    string
	Password    string
	Server      string
	Interactive bool
	Verbose     bool
	Timeout     time.Duration
}

func ConfigInit() *Config {
	return &Config{
		Server:  "https://oidc-server.example.com",
		Timeout: 60 * time.Second,
	}
}
