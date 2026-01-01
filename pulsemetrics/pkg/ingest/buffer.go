package ingest

import (
	"log"
	"sync/atomic"
)

// Buffer holds incoming metric events
type Buffer struct {
	ch chan MetricEvent
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		ch: make(chan MetricEvent, size),
	}
}

func (b *Buffer) Push(event MetricEvent) {
	select {
	case b.ch <- event:
		atomic.AddUint64(&ingestionCount, 1)
	default:
		
		atomic.AddUint64(&metricsDropped, 1)
		log.Println("buffer full, dropping metric")
	}
}

// Input exposes the underlying channel for consumers.
func (b *Buffer) Input() <-chan MetricEvent { return b.ch }

// Close closes the buffer channel to signal consumers.
func (b *Buffer) Close() { close(b.ch) }
