package storage

import (
	"log"

	"github.com/vinayak55jain/pulsemetrics/pkg/ingest"
)

type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Save(batch []ingest.MetricEvent) error {
	log.Printf("saved batch of %d metrics\n", len(batch))
	return nil
}
