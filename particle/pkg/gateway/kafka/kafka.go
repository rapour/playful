package kafka

import (
	"context"

	"playful/particles/pkg/domain"
	"playful/particles/tools/kafka"

	"github.com/rs/zerolog/log"
	kafka_go "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaClient interface {
	SendLocarion(ctx context.Context, p domain.Location) error
	Shutdown() error
}

type kafkaClient struct {
	config kafka.Config
	writer kafka_go.Writer
}

func (k *kafkaClient) Shutdown() error {

	return k.writer.Close()

}

func (k *kafkaClient) SendLocarion(ctx context.Context, l domain.Location) error {

	loc := Location{
		Id:        1,
		Longitude: l.Longitude,
		Altitude:  l.Altitude,
		Timestamp: l.Timestamp,
	}

	locationBytes, err := proto.Marshal(&loc)
	if err != nil {
		return err
	}

	err = k.writer.WriteMessages(context.TODO(), kafka_go.Message{
		//Key:   []byte(loc.Id),
		Value: locationBytes,
	})

	return err
}

func NewKafkaClient(c kafka.Config) KafkaClient {

	dialer := c.Dialer()

	sharedTransport := &kafka_go.Transport{
		SASL: dialer.SASLMechanism,
		TLS:  dialer.TLS,
	}

	w := kafka_go.Writer{
		Addr:                   kafka_go.TCP(c.BootstrapServer),
		Topic:                  c.Topic,
		AllowAutoTopicCreation: true,
		Logger:                 kafka_go.LoggerFunc(log.Debug().Msgf),
		ErrorLogger:            kafka_go.LoggerFunc(log.Error().Msgf),
		Balancer:               &kafka_go.LeastBytes{},
		Transport:              sharedTransport,
	}

	return &kafkaClient{
		config: c,
		writer: w,
	}
}
