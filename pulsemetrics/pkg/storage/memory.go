package storage

import (
	"log"

	m "github.com/vinayak55jain/pulsemetrics/pkg/model"
)

type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (ms *MemoryStore) Save(batch []m.MetricEvent) error {
	log.Printf("saved batch of %d metrics\n", len(batch))
	return nil
}
