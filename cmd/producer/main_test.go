package main

import (
	"testing"

	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/config"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/generator"
)

func BenchmarkGenerateTransaction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generator.GenerateTransaction()
	}
}

func LoadTestConfig() *config.Config {
	return &config.Config{
		KafkaBroker: "localhost:9092",
		KafkaTopic:  "smart-grid-readings-test", // Use a separate test topic
	}
}
