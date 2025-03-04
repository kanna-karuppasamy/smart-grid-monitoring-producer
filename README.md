# Energy Monitoring Kafka Producer

This project simulates an energy monitoring system that generates smart meter data and publishes it to Kafka. It's designed to demonstrate high-throughput data ingestion for real-time analytics.

## Features

- Generates realistic energy consumption data from simulated smart meters
- Supports different building types (residential, commercial, industrial)
- Includes geographic information for spatial analysis
- Simulates meter faults and offline status for monitoring
- Efficiently publishes to Kafka with batching and compression

## Use Cases

The generated data supports the following analytics use cases:

- **Live Energy Consumption Per Second:** Track real-time energy usage
- **Power Usage Trends (Last 24 Hours):** Analyze consumption patterns
- **List of Smart Meters with Faults or Offline Status:** Monitor system health
- **Map of High-Energy Consumption Regions:** Identify consumption hotspots

## Project Structure

```
smart-grid-monitoring-producer/
├── cmd/
│   └── producer/
│       └── main.go        // Entry point for the application
├── internal/
│   ├── config/
│   │   └── config.go      // Configuration management
│   ├── generator/
│   │   └── generator.go   // Mock data generation
│   ├── kafka/
│   │   └── producer.go    // Kafka publishing
│   └── models/
│       └── transaction.go // Data models
```

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Kafka cluster (or local Kafka setup)

### Installation

1. Clone the repository
```bash
git clone https://github.com/kanna-karuppasamy/smart-grid-monitoring-producer.git
cd smart-grid-monitoring-producer
```

2. Install dependencies
```bash
go mod download
```

### Configuration

Edit the configuration in `config.json` or use environment variables:

```json
{
  "kafka": {
    "brokers": ["localhost:9092"],
    "topic": "energy-transactions",
    "batchSize": 100,
    "compression": "snappy"
  },
  "generator": {
    "totalTransactions": 1000000,
    "meterCount": 10,
    "faultProbability": 0.01,
    "offlineProbability": 0.005
  }
}
```

### Running

```bash
go run cmd/producer/main.go
```

Or build and run the binary:

```bash
go build -o energy-producer cmd/producer/main.go
./energy-producer
```

## Performance

The application is designed for high throughput:

- Batched message publishing
- Message compression
- Configurable transaction rate
- Progress monitoring

## Next Steps

1. Create a consumer application to process the data
2. Implement a real-time dashboard using the data
3. Add more sophisticated data generation patterns
4. Integrate with a time-series database for historical analysis

## License

MIT