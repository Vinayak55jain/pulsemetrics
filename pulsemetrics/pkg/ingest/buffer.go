package ingest

import "log"

// Buffer holds incoming metric events
type Buffer struct {
	ch chan MetricEvent
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		ch: make(chan MetricEvent, size),
	}
}

// Push is NON-BLOCKING
func (b *Buffer) Push(event MetricEvent) {
	select {
	case b.ch <- event:
		// accepted
	default:
		// buffer full → drop
		log.Println("buffer full, dropping metric")
	}
}
