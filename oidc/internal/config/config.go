package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	CronConfig string
}

type CronConfig struct {
	Cron string
}

func NewConfig() (*Config, error) {
	log.Println("Config init started")
	c := &Config{
		CronConfig: "",
	}
	err := c.CronInit()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) CronInit() error {
	var exists bool
	c.CronConfig, exists = os.LookupEnv("CRON")
	if !exists {
		return fmt.Errorf("")
	}
	return nil
}
