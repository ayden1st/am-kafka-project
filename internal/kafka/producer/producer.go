package producer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"

	"am-kafka-project/pkg/common"
)

type ProducerService struct {
	Producer       *kafka.Producer
	Topic          string
	SchemaRegistry schemaregistry.Client
	Serializer     *jsonschema.Serializer
}

func NewProducerService() *ProducerService {
	return &ProducerService{}
}

// TODO Add configure parameters
func (p *ProducerService) Configure() (err error) {
	p.Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": common.GetEnv("AKP_KAFKA_BROKERS", "localhost:9092"),
		"client.id":         common.GetEnv("AKP_KAFKA_CLIENT_ID", "am-alerts-dev"),
		"acks":              "all",
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %s", err)
	}

	p.SchemaRegistry, err = schemaregistry.NewClient(schemaregistry.NewConfig(common.GetEnv("AKP_SCHEMA_REGISTRY", "http://localhost:8081")))
	if err != nil {
		return fmt.Errorf("failed to create schema registry client: %s", err)
	}

	p.Serializer, err = jsonschema.NewSerializer(p.SchemaRegistry, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		return fmt.Errorf("failed to create schema registry serializer: %s", err)
	}
	p.Topic = common.GetEnv("AKP_KAFKA_TOPIC", "alerts")

	return nil
}

func (p *ProducerService) Close() {
	p.Producer.Close()
	p.SchemaRegistry.Close()
}

func (p *ProducerService) Push(topic string, key, value []byte) (err error) {
	message := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
		Timestamp:      time.Now(),
	}

	return p.Producer.Produce(&message, nil)
}
