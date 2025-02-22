# Smart Grid Monitoring - Kafka Producer

## Overview
The `smart-grid-monitoring-producer` is a **Kafka producer** that simulates smart meters sending real-time energy consumption data to a Kafka topic. This microservice generates high-frequency messages to support real-time monitoring and analytics in the Smart Grid Monitoring system, with a goal of producing **1 million transactions per second (1M TPS)**.

## Features
- Generates **real-time energy consumption** data from simulated smart meters.
- Publishes messages to **Kafka topic (`smart-grid-data`)** at extremely high throughput.
- Configurable message frequency and payload structure.
- Supports **manual scaling** by running multiple instances.
- Optimized for **high-throughput performance** without Kubernetes orchestration.

---

## Repository Structure
```plaintext
smart-grid-monitoring-producer/
│── config/
│   ├── config.yaml            # Configuration file (Kafka brokers, topic, etc.)
│── internal/
│   ├── generator/             # Logic for generating smart meter data
│   │   ├── generator.go
│   ├── producer/              # Kafka producer logic
│   │   ├── producer.go
│── scripts/
│   ├── start.sh               # Shell script to start the producer
│── Dockerfile                 # Dockerfile for containerization (optional)
│── go.mod                      # Go module dependencies
│── go.sum                      # Go module checksum
│── main.go                     # Main entry point of the producer
│── README.md                   # Documentation
```

---

## Setup Instructions

### 1. Prerequisites
- **Golang** (>=1.19)
- **Kafka** (Running instance with optimized partitions)
- **High-performance infrastructure** (Multiple Kafka brokers, optimized networking)

### 2. Installation
Clone the repository:
```bash
git clone https://github.com/your-org/smart-grid-monitoring-producer.git
cd smart-grid-monitoring-producer
```

### 3. Configuration
Edit `config/config.yaml` to update Kafka settings:
```yaml
kafka:
  brokers: ["localhost:9092"]
  topic: "smart-grid-data"
  client_id: "smart-grid-producer"
message_rate: 1000000 # Messages per second
```

### 4. Running Locally
Run the producer locally:
```bash
go run main.go
```

Or build and execute:
```bash
go build -o producer
./producer
```

### 5. Docker Build & Run (Optional)
```bash
docker build -t smart-grid-producer .
docker run --rm -e KAFKA_BROKER=localhost:9092 smart-grid-producer
```

---

## Message Format
Each message contains:
```json
{
  "meter_id": "12345",
  "timestamp": "2025-02-22T12:34:56Z",
  "energy_consumption": 2.5,
  "status": "Online"
}
```

---

## Scaling & Performance
- **Goal: 1M TPS**
- Run **multiple producer instances** to increase throughput.
- Optimize **Kafka partitions** and producer batch size.
- Use **Kafka compression (Snappy, LZ4)** to reduce network overhead.
- Deploy on **high-performance servers** (optimized CPU/memory).

---

## Monitoring & Logs
- Logs can be viewed using `tail -f logs/output.log`.
- Kafka metrics can be monitored using **Kafka Exporter**.
- Utilize **Jaeger for tracing high-throughput message flows**.

---

## License
MIT License
