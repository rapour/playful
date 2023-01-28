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
	hiddenCount := 0
	var inertia float64 = 0.2
	for {

		loc := domain.Location{Timestamp: int32(time.Now().Unix())}

		if hiddenCount == 0 {
			loc.Altitude = int32(yy)
			loc.Longitude = int32(xx)

			coin := rand.Intn(100)
			if coin > 95 {
				hiddenCount = 5
			}
		} else {
			hiddenCount--
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

		if math.Abs(vx) > 5 {
			vx = 5 * (vx / math.Abs(vx))
		}
		if math.Abs(vy) > 5 {
			vy = 5 * (vy / math.Abs(vy))
		}

	}

	log.Printf("Particle is dying")
}
