package main

import (
	"math/rand"
	"os"
	"time"

	controller "playful/app/pkg/controller/kafka"
	"playful/app/tools/kafka"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	rand.Seed(time.Now().Unix())

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("application started")

	kafkaConfig, err := kafka.NewKafkaConfigFromENV("")
	if err != nil {
		log.Panic().Msgf("error obtaining kafka config: %v", err)
	}

	kafkaContoller := controller.NewKafkaController(kafkaConfig)

	errChan := kafkaContoller.Manager()

	err = <-errChan
	log.Err(err).Msg("kafka controller failed")

}
