package storage

import "github.com/vinayak55jain/pulsemetrics/pkg/ingest"

type Storage interface {
	Save(batch []ingest.MetricEvent) error
}
