package storage

import m "github.com/vinayak55jain/pulsemetrics/pkg/model"

// Storage is the minimal persistence/aggregation contract used by the
// ingestion worker pool. It depends only on the shared model package to avoid
// import cycles with the ingest package.
type Storage interface {
	Save(batch []m.MetricEvent) error
}
