package cassandra

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address           string `envconfig:"cassandra_listen_address" required:"true"`
	Keyspace          string `envconfig:"cassandra_keyspace" default:"playful"`
	ReplicationFactor int    `envconfig:"cassandra_replication_factor" default:"3"`
}

func NewCassandraConfigFromEnv(prefix string) (Config, error) {

	var config Config
	err := envconfig.Process(prefix, &config)
	if err != nil {
		return config, fmt.Errorf("error obtaining cassandra config: %v", err)
	}

	return config, nil
}
