package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	Kafka     KafkaConfig     `json:"kafka"`
	Generator GeneratorConfig `json:"generator"`
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Brokers         []string `json:"brokers"`
	Topic           string   `json:"topic"`
	BatchSize       int      `json:"batchSize"`
	Compression     string   `json:"compression"`
	RequiredAcks    int      `json:"requiredAcks"`
	LingerMs        int      `json:"lingerMs"`
	MaxMessageBytes int      `json:"maxMessageBytes"`
}

// GeneratorConfig holds data generation configuration
type GeneratorConfig struct {
	TotalTransactions  int      `json:"totalTransactions"`
	MeterCount         int      `json:"meterCount"`
	FaultProbability   float64  `json:"faultProbability"`
	OfflineProbability float64  `json:"offlineProbability"`
	Regions            []Region `json:"regions"`
}

// Region represents a geographic region for data generation
type Region struct {
	Name      string  `json:"name"`
	MinLat    float64 `json:"minLat"`
	MaxLat    float64 `json:"maxLat"`
	MinLong   float64 `json:"minLong"`
	MaxLong   float64 `json:"maxLong"`
	MeterPerc float64 `json:"meterPercentage"` // Percentage of meters in this region
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Kafka: KafkaConfig{
			Brokers:         []string{"localhost:9092"},
			Topic:           "smart-grid-readings",
			BatchSize:       5000,    // Increase batch size
			Compression:     "lz4",   // Zstd is faster & compresses better than lz4
			RequiredAcks:    1,       // Use 1 instead of 0 to avoid excessive retries
			MaxMessageBytes: 2000000, // Increase max message size
			LingerMs:        20,      // Increase linger time for better batching
		},
		Generator: GeneratorConfig{
			TotalTransactions:  0,
			MeterCount:         100,
			FaultProbability:   0.01,
			OfflineProbability: 0.005,
			Regions: []Region{
				{
					Name:      "Urban",
					MinLat:    40.7128,
					MaxLat:    40.8128,
					MinLong:   -74.0060,
					MaxLong:   -73.9060,
					MeterPerc: 0.6,
				},
				{
					Name:      "Suburban",
					MinLat:    40.6128,
					MaxLat:    40.7128,
					MinLong:   -74.1060,
					MaxLong:   -74.0060,
					MeterPerc: 0.3,
				},
				{
					Name:      "Rural",
					MinLat:    40.5128,
					MaxLat:    40.6128,
					MinLong:   -74.2060,
					MaxLong:   -74.1060,
					MeterPerc: 0.1,
				},
			},
		},
	}
}

// Load loads configuration from a file or environment variables
// Falls back to defaults if no configuration is found
func Load() (*Config, error) {
	config := DefaultConfig()

	// Try to load from config file if it exists
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	// Override with environment variables if they exist
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		var brokerList []string
		if err := json.Unmarshal([]byte(brokers), &brokerList); err == nil {
			config.Kafka.Brokers = brokerList
		}
	}

	if topic := os.Getenv("KAFKA_TOPIC"); topic != "" {
		config.Kafka.Topic = topic
	}

	return config, nil
}
