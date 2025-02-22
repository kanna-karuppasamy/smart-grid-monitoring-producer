package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/config"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/generator"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/kafka"
	models "github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/models"
)

func main() {
	cfg := config.LoadConfig()

	// Create Kafka producer
	producer, err := kafka.NewProducer(cfg.KafkaBroker)
	if err != nil {
		fmt.Printf("Failed to create Kafka producer: %v\n", err)
		return
	}
	defer producer.Close()

	// Configuration
	const (
		batchSize  = 1_000 // Reduced batch size for Kafka
		numWorkers = 2     // Number of concurrent workers
	)

	// Create channels
	transactionChan := make(chan models.Transaction, batchSize)
	doneChan := make(chan bool)

	// Start generator goroutines
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				select {
				case <-doneChan:
					return
				default:
					transactionChan <- generator.GenerateTransaction()
				}
			}
		}()
	}

	// Process and send transactions
	var wg sync.WaitGroup
	transactions := make([]models.Transaction, 0, batchSize)
	counter := 0
	startTime := time.Now()

	fmt.Println("Starting data generation and sending to Kafka...")

	for {
		transaction := <-transactionChan
		transactions = append(transactions, transaction)
		counter++

		// Send batch to Kafka when batch is full
		if len(transactions) >= batchSize {
			wg.Add(1)
			batch := make([]models.Transaction, len(transactions))
			copy(batch, transactions)
			go kafka.SendToKafka(producer, batch, cfg.KafkaTopic, &wg)
			transactions = transactions[:0]

			// Progress update
			elapsed := time.Since(startTime).Seconds()
			tps := float64(counter) / elapsed
			fmt.Printf("Transactions generated: %d (%.2f TPS)\n", counter, tps)
		}
	}
}
