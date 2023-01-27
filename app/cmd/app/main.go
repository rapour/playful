package main

import (
	"math/rand"
	"os"
	"time"

	controller "playful/app/pkg/controller/kafka"
	cassandra_repository "playful/app/pkg/repository/cassandra"
	"playful/app/pkg/service/playful"
	"playful/app/tools/cassandra"
	"playful/app/tools/kafka"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	rand.Seed(time.Now().Unix())

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("application started")

	cassandraConfig, err := cassandra.NewCassandraConfigFromEnv("")
	if err != nil {
		log.Panic().Msgf("error obtaining cassandra config: %v", err)
	}

	cassandraClient, err := cassandra.NewCassandraClient(cassandraConfig)
	if err != nil {
		log.Panic().Msgf("error connecting to cassandra: %v", err)
	}
	defer cassandraClient.Close()

	columnRepository := cassandra_repository.NewColumnRepository(cassandraClient)

	playfulService, err := playful.NewPlayfulService(columnRepository)
	if err != nil {
		log.Panic().Msgf("error initiating the service: %v", err)
	}

	kafkaConfig, err := kafka.NewKafkaConfigFromENV("")
	if err != nil {
		log.Panic().Msgf("error obtaining kafka config: %v", err)
	}

	kafkaContoller := controller.NewKafkaController(kafkaConfig, playfulService)

	errChan := kafkaContoller.Manager()

	err = <-errChan
	log.Err(err).Msg("kafka controller failed")

}
