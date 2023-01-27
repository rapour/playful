package http

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port    string `envconfig:"http_server_port" required:"true"`
	Address string `ignore:"true"`
}

func NewHttpConfig(prefix string) (hc Config, err error) {

	var conf Config
	if err := envconfig.Process(prefix, &conf); err != nil {
		return Config{}, fmt.Errorf("error obtaining http configuration: %w", err)
	}

	conf.Address = "localhost"

	return conf, nil
}
