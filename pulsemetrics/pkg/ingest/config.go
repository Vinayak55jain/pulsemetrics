package ingest

import "time"

type Config struct {
	BatchSize     int
	FlushInterval time.Duration
	BufferSize    int
	WorkerCount   int
}

// DefaultConfig gives sane production defaults
func DefaultConfig() Config {
	return Config{
		BatchSize:     10,
		FlushInterval: 10 * time.Millisecond,
		BufferSize:    1000,
		WorkerCount:   4,
	}
}
