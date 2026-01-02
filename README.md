# 🚀 PulseMetrics

**PulseMetrics** is a **high-throughput, buffered metrics ingestion system** written in Go.
It is designed to handle **millions of small, high-frequency events** (metrics, logs, telemetry) without overwhelming backend services.

The project demonstrates **real-world system design patterns** used in observability, analytics, and event-driven platforms.

---

## 📌 Why PulseMetrics?

In large-scale systems, metrics and telemetry are:

* High frequency
* Small payload
* Eventually consistent

Sending every event synchronously to a database or message broker causes:

* High latency
* Backend overload
* Poor scalability
* Increased infrastructure cost

**PulseMetrics solves this by decoupling ingestion from processing.**

---

## 🧠 Core Design Principles

* **Non-blocking HTTP ingestion**
* **Bounded in-memory buffering with backpressure**
* **Micro-batching (size-based OR time-based flush)**
* **Asynchronous worker pools**
* **Graceful shutdown with no data loss**
* **Pluggable storage backends (Memory / Kafka / DB)**

---

## 🏗️ High-Level Architecture

```
Client
  ↓
HTTP Ingestion Layer (202 Accepted)
  ↓
In-Memory Buffer (bounded)
  ↓
Batcher (BatchSize OR FlushInterval)
  ↓
Worker Pool (async processing)
  ↓
Storage Layer (Kafka / DB / Memory)
```

This architecture mirrors patterns used in **production telemetry pipelines**.

---

## 📂 Project Structure

```
pulsemetrics/
├── cmd/
│   └── server/          # Standalone ingestion service
├── pkg/
│   ├── api/             # HTTP ingestion layer
│   ├── ingest/          # Buffer, batcher, workers, ingestor
│   └── storage/         # Storage interfaces & implementations
├── go.mod
└── README.md
```

---

## 🌐 HTTP API

### `POST /ingest`

Accepts a single metric/event.
The request is **accepted immediately** and processed asynchronously.

#### Request Body

```json
{
  "type": "page_view",
  "name": "/home",
  "value": 1,
  "tags": {
    "browser": "chrome",
    "user": "anon"
  }
}
```

#### Response

```http
202 Accepted
```

> The server does **not wait** for batching or storage.

---

## ⚙️ Configuration

Key configuration parameters:

| Parameter       | Description                      |
| --------------- | -------------------------------- |
| `BufferSize`    | Max events stored in memory      |
| `BatchSize`     | Max events per batch             |
| `FlushInterval` | Max wait time before batch flush |
| `WorkerCount`   | Number of concurrent workers     |

Example defaults:

```go
BufferSize    = 1000
BatchSize     = 10
FlushInterval = 10 * time.Millisecond
WorkerCount   = 4
```

---

## 🧵 Backpressure & Reliability

* Buffer is **bounded**
* On overflow, **metrics are dropped (not requests)**
* Server remains responsive during traffic spikes
* Workers isolate slow storage from ingestion

This ensures **graceful degradation instead of system failure**.

---

## 🧰 Using PulseMetrics as a Library

PulseMetrics can be embedded into any Go service.

```go
store := storage.NewMemoryStore()
cfg := ingest.DefaultConfig()

ingestor, _ := ingest.NewIngestor(cfg, store)
defer ingestor.Close()

ingestor.Push(metricEvent)
```

---

## 🔌 Using PulseMetrics with Kafka

PulseMetrics is designed to integrate naturally with **Kafka**.

### Why Kafka?

* Durable event storage
* Horizontal scalability
* Fan-out consumption
* Event replay

### Kafka Storage Implementation (Example)

```go
type KafkaStore struct {
    producer *kafka.Producer
    topic    string
}

func (k *KafkaStore) Save(batch []ingest.MetricEvent) error {
    for _, event := range batch {
        data, _ := json.Marshal(event)
        k.producer.Produce(&kafka.Message{
            TopicPartition: kafka.TopicPartition{
                Topic: &k.topic,
                Partition: kafka.PartitionAny,
            },
            Value: data,
        }, nil)
    }
    return nil
}
```

### Kafka Flow

```
PulseMetrics → Kafka Topic
                    ↓
      Consumers (Analytics / Storage / Alerts)
```

This enables **independent scaling and downstream processing**.

---

## 🧯 Graceful Shutdown

On shutdown:

1. HTTP server stops accepting requests
2. Batcher flushes remaining events
3. Batch channel is closed
4. Workers finish in-flight batches
5. Process exits cleanly

No partial writes. No goroutine leaks.

---

## 📈 Where This Pattern Is Used

* Observability & monitoring systems
* Analytics pipelines
* Log aggregation platforms
* Telemetry agents
* High-throughput backend services

---

## 📐 System Design

This repository includes:

* High-level architecture diagrams
* Detailed flow diagrams
* Failure & backpressure flows
* Shutdown sequence

All designs are **production-inspired**.

---
## 🔗 Repository

**GitHub:**
https://github.com/Vinayak55jain/pulsemetrics.git

## 🧑‍💻 Purpose

Built to learn and demonstrate:

* Backend system design
* Go concurrency patterns
* High-throughput ingestion
* Production-grade architecture decisions
