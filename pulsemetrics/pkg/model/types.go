package model

import "time"

// MetricEvent is the canonical data model for incoming metrics.
type MetricEvent struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}
