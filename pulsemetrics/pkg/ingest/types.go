package ingest

import "time"

// MetricEvent represents ONE metric sent by client
type MetricEvent struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
	Timestamp time.Time         `json:"timestamp"`
}
