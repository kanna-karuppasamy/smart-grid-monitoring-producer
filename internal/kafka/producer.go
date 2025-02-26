package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/config"
)

// Producer is responsible for publishing messages to Kafka
type Producer struct {
	producer sarama.SyncProducer
	config   config.KafkaConfig
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	// Create Kafka producer configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true

	// Enable compression if specified
	switch cfg.Compression {
	case "gzip":
		kafkaConfig.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		kafkaConfig.Producer.Compression = sarama.CompressionSnappy
	case "lz4":
		kafkaConfig.Producer.Compression = sarama.CompressionLZ4
	case "zstd":
		kafkaConfig.Producer.Compression = sarama.CompressionZSTD
	default:
		kafkaConfig.Producer.Compression = sarama.CompressionNone
	}

	kafkaConfig.Producer.Flush.Frequency = time.Duration(cfg.LingerMs)
	kafkaConfig.Producer.RequiredAcks = sarama.NoResponse

	// Create producer
	producer, err := sarama.NewSyncProducer(cfg.Brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		config:   cfg,
	}, nil
}

// Publish publishes a single message to Kafka
func (p *Producer) Publish(ctx context.Context, topic string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(jsonData),
		Timestamp: time.Now(),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}

// PublishBatch publishes a batch of messages to Kafka
func (p *Producer) PublishBatch(ctx context.Context, topic string, messages []interface{}) error {
	batch := make([]*sarama.ProducerMessage, len(messages))

	for i, message := range messages {
		jsonData, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		batch[i] = &sarama.ProducerMessage{
			Topic:     topic,
			Value:     sarama.ByteEncoder(jsonData),
			Timestamp: time.Now(),
		}
	}

	return p.producer.SendMessages(batch)
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
