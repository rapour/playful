package kafka

import (
	"context"

	"playful/app/pkg/controller"
	"playful/app/tools/kafka"

	"github.com/rs/zerolog/log"
	kafka_go "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type kafkaController struct {
	config       kafka.Config
	errChan      chan error
	innerErrChan chan error
}

func (k *kafkaController) Worker() {

	defer func() {

		if r := recover(); r != nil {
			log.Error().Msgf("kafka worker recovering from panic: %v", r)
		}
		go k.Worker()
	}()

	dialer := k.config.Dialer()

	reader := kafka_go.NewReader(kafka_go.ReaderConfig{
		Brokers:     []string{k.config.BootstrapServer},
		Topic:       k.config.Topic,
		GroupID:     "workers",
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		Logger:      kafka_go.LoggerFunc(log.Debug().Msgf),
		ErrorLogger: kafka_go.LoggerFunc(log.Error().Msgf),
		Dialer:      &dialer,
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.TODO())
		if err != nil {
			k.innerErrChan <- err
			return
		}
		log.Info().Msgf("message received at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var message Location
		err = proto.Unmarshal(m.Value, &message)
		if err != nil {
			k.innerErrChan <- err
			return
		}

		// process the message

		// err = reader.CommitMessages(context.TODO(), m)
		// if err != nil {
		// 	k.innerErrChan <- err
		// 	return
		// }
	}

}

func (k *kafkaController) Manager() chan error {

	WorkerNumber := 1
	go func() {

		for i := 0; i < WorkerNumber; i++ {
			go k.Worker()
		}

		// communicate inner channel to the exposed error channel if necessary
		for err := range k.innerErrChan {
			log.Error().Msgf("error in kafka worker: %s", err)
		}

	}()

	return k.errChan
}

func NewKafkaController(c kafka.Config) controller.KafkaController {

	return &kafkaController{
		config:  c,
		errChan: make(chan error),
	}
}
