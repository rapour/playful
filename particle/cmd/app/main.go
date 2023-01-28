package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"playful/particles/pkg/domain"
	gateway "playful/particles/pkg/gateway/kafka"
	"playful/particles/tools/kafka"
)

func main() {

	rand.Seed(time.Now().Unix())

	log.Printf("Particle is starting to play")

	kafkaConfig, err := kafka.NewKafkaConfigFromENV("")
	if err != nil {
		log.Panicf("error obtaining kafka config: %v", err)
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

	log.Printf("Particle is dying")
}
