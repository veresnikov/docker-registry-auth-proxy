package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

func parseEnv() (*Config, error) {
	c := new(Config)
	if err := envconfig.Process(applicationID, c); err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	ServeHTTPAddr    string        `envconfig:"serve_http_addr" default:":8000"`
	ServeWriteTimout time.Duration `envconfig:"serve_write_timeout" default:"1h"`
	ServeReadTimout  time.Duration `envconfig:"serve_read_timeout" default:"1h"`
	ServeGracePeriod time.Duration `envconfig:"serve_grace_period" default:"30s"`

	AccessConfigPath string `envconfig:"access_config_path" default:"/app/.authproxy/config.json"`
	RegistryAddress  string `envconfig:"registry_address" required:"true"`
}
