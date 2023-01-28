package main

import (
	"log"
	"math/rand"
	"time"

	http_controller "playful/app/pkg/controller/http"
	kafka_controller "playful/app/pkg/controller/kafka"
	cassandra_repository "playful/app/pkg/repository/cassandra"
	"playful/app/pkg/service/playful"
	"playful/app/tools/cassandra"
	"playful/app/tools/http"
	"playful/app/tools/kafka"
)

func main() {

	rand.Seed(time.Now().Unix())

	httpConfig, err := http.NewHttpConfig("")
	if err != nil {
		log.Printf("error obtaining http config: %v", err)
	}

	cassandraConfig, err := cassandra.NewCassandraConfigFromEnv("")
	if err != nil {
		log.Panicf("error obtaining cassandra config: %v", err)
	}

	cassandraClient, err := cassandra.NewCassandraClient(cassandraConfig)
	if err != nil {
		log.Panicf("error connecting to cassandra: %v", err)
	}
	defer cassandraClient.Close()

	columnRepository := cassandra_repository.NewColumnRepository(cassandraClient)

	playfulService, err := playful.NewPlayfulService(columnRepository)
	if err != nil {
		log.Panicf("error initiating the service: %v", err)
	}

	kafkaConfig, err := kafka.NewKafkaConfigFromENV("")
	if err != nil {
		log.Panicf("error obtaining kafka config: %v", err)
	}

	kafkaContoller := kafka_controller.NewKafkaController(kafkaConfig, playfulService)

	httpController := http_controller.NewHttpController(httpConfig, playfulService)

	kafkaErrChan := kafkaContoller.Manager()
	httpErrChan := httpController.Manager()

	select {
	case err := <-kafkaErrChan:
		log.Panicf("kafka controller failed: %v", err)

	case err := <-httpErrChan:
		log.Panicf("http controller failed: %v", err)
	}

}
