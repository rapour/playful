package main

import (
	"context"
	"log"
	"math"
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

	var vx float64 = 4
	var vy float64 = 4
	var xx float64 = float64(rand.Intn(100))
	var yy float64 = float64(rand.Intn(100))
	var inertia float64 = 0.5
	for {

		loc := domain.Location{
			Longitude: int32(xx),
			Altitude:  int32(yy),
			Timestamp: int32(time.Now().Unix()),
		}

		err := kafkaClient.SendLocarion(context.TODO(), loc)
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(15 * time.Millisecond)

		// calculate new location
		xx = xx + vx
		yy = yy + vy

		if xx >= 200 || xx <= 0 {
			xx = math.Min(200, math.Max(xx, 0))
			vx = -vx
		}

		if yy >= 200 || yy <= 0 {
			yy = math.Min(200, math.Max(yy, 0))
			vy = -vy
		}

		// calcualte new velocity
		vx = vx + vy*rand.NormFloat64()*inertia
		vy = vy + vx*rand.NormFloat64()*inertia

		if vx > 100 {
			vx = 100
		}
		if vy > 100 {
			vy = 100
		}

	}

	log.Printf("Particle is dying")
}
