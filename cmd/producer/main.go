package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/config"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/generator"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/kafka"
)

func main() {
	log.SetOutput(os.Stdout)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Create data generator
	gen := generator.NewGenerator(cfg.Generator)

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, initiating shutdown", sig)
		cancel()
	}()

	// Start the production process
	var wg sync.WaitGroup

	startTime := time.Now()
	var totalMessages int64 = 0
	var continuousMode bool = cfg.Generator.TotalTransactions <= 0

	// Use a channel for pre-generating messages
	messageChannel := make(chan interface{}, 50000) // Buffer 50K messages

	// Pre-generate messages in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(messageChannel)

		// If TotalTransactions <= 0, run continuously until context is cancelled
		if continuousMode {
			log.Println("Starting continuous transaction generation...")
			for ctx.Err() == nil {
				select {
				case messageChannel <- gen.GenerateTransaction():
					// Message added to channel
				case <-ctx.Done():
					log.Println("Stopping message generation due to context cancellation")
					return
				}
			}
		} else {
			log.Printf("Starting to generate %d transactions...", cfg.Generator.TotalTransactions)
			for i := 0; i < cfg.Generator.TotalTransactions && ctx.Err() == nil; i++ {
				select {
				case messageChannel <- gen.GenerateTransaction():
					// Message added to channel
				case <-ctx.Done():
					log.Println("Stopping message generation due to context cancellation")
					return
				}
			}
		}
	}()

	// Determine optimal number of producer goroutines
	numProducers := runtime.NumCPU()
	log.Printf("Starting %d producer goroutines", numProducers)

	// Launch multiple producer goroutines
	for i := 0; i < numProducers; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()

			batchSize := cfg.Kafka.BatchSize
			batch := make([]interface{}, 0, batchSize)
			localCounter := 0
			lastLogTime := time.Now()

			for transaction := range messageChannel {
				if ctx.Err() != nil {
					break
				}

				// Add to batch
				batch = append(batch, transaction)

				// When batch is full, send to Kafka
				if len(batch) >= batchSize {
					if err := producer.PublishBatch(ctx, cfg.Kafka.Topic, batch); err != nil {
						log.Printf("Producer %d: Failed to publish batch: %v", producerID, err)
					} else {
						newCount := atomic.AddInt64(&totalMessages, int64(len(batch)))
						localCounter += len(batch)

						// In continuous mode, log based on time instead of message count
						if continuousMode {
							now := time.Now()
							if now.Sub(lastLogTime) >= 5*time.Second {
								elapsed := time.Since(startTime)
								rate := float64(newCount) / elapsed.Seconds()
								log.Printf("Progress: %d transactions published (%.2f msgs/sec)",
									newCount, rate)
								lastLogTime = now
							}
						} else if newCount%100000 == 0 {
							elapsed := time.Since(startTime)
							rate := float64(newCount) / elapsed.Seconds()
							log.Printf("Progress: %d/%d transactions published (%.2f msgs/sec)",
								newCount, cfg.Generator.TotalTransactions, rate)
						}
					}
					// Reuse the slice without reallocating
					batch = batch[:0]
				}
			}

			// Send any remaining transactions in the batch
			if len(batch) > 0 && ctx.Err() == nil {
				if err := producer.PublishBatch(ctx, cfg.Kafka.Topic, batch); err != nil {
					log.Printf("Producer %d: Failed to publish final batch: %v", producerID, err)
				} else {
					atomic.AddInt64(&totalMessages, int64(len(batch)))
				}
			}

			log.Printf("Producer %d completed, published %d messages", producerID, localCounter)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print final statistics
	elapsed := time.Since(startTime)
	finalTotal := atomic.LoadInt64(&totalMessages)
	rate := float64(finalTotal) / elapsed.Seconds()
	log.Printf("Completed! Published %d transactions in %.2f seconds (%.2f msgs/sec)",
		finalTotal, elapsed.Seconds(), rate)
}
