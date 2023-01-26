package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"playful/particles/pkg/domain"
	gateway "playful/particles/pkg/gateway/kafka"
	"playful/particles/tools/kafka"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	rand.Seed(time.Now().Unix())

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Particle is starting to play")

	kafkaConfig, err := kafka.NewKafkaConfigFromENV("")
	if err != nil {
		log.Panic().Msgf("error obtaining kafka config: %v", err)
	}

	kafkaClient := gateway.NewKafkaClient(kafkaConfig)
	defer kafkaClient.Shutdown()

	for {
		loc := domain.Location{
			Longitude: int32(rand.Intn(100)),
			Altitude:  int32(rand.Intn(100)),
			Timestamp: int32(time.Now().Unix()),
		}

		kafkaClient.SendLocarion(context.TODO(), loc)
		time.Sleep(time.Second)
	}

	log.Info().Msg("Particle is dying")
}
