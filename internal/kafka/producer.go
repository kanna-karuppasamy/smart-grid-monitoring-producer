package kafka

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	models "github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/models"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(broker string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Flush.Messages = 1000                  // Try increasing batch size
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) Close() {
	p.producer.Close()
}

func SendToKafka(producer *Producer, transactions []models.Transaction, topic string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, transaction := range transactions {
		jsonData, err := json.Marshal(transaction)
		if err != nil {
			fmt.Printf("Error marshaling transaction: %v\n", err)
			continue
		}

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(transaction.MeterID),
			Value: sarama.ByteEncoder(jsonData),
		}

		_, _, err = producer.producer.SendMessage(msg)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			continue
		}
	}
}
