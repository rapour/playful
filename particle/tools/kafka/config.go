package kafka

import (
	"github.com/kelseyhightower/envconfig"

	"time"

	kafka_go "github.com/segmentio/kafka-go"
)

type Config struct {
	BootstrapServer string `envconfig:"kafka_bootstrap_server"`
	Topic           string `envconfig:"kafka_topic"`
	dialer          kafka_go.Dialer
}

func (c Config) Dialer() kafka_go.Dialer {
	return c.dialer
}

func NewKafkaConfigFromENV(prefix string) (Config, error) {

	var config Config
	err := envconfig.Process(prefix, &config)

	d := kafka_go.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       nil,
	}
	config.dialer = d

	return config, err
}
