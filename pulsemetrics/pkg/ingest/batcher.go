package ingest

import "time"

type Batcher struct {
	buffer *Buffer
	cfg    Config
	out    chan []MetricEvent
	stop   chan struct{}
}

func NewBatcher(buffer *Buffer, cfg Config) *Batcher {
	return &Batcher{
		buffer: buffer,
		cfg:    cfg,
		out:    make(chan []MetricEvent),
		stop:   make(chan struct{}),
	}
}

func (b *Batcher) Start() {
	go func() {
		ticker := time.NewTicker(b.cfg.FlushInterval)
		defer ticker.Stop()

		batch := make([]MetricEvent, 0, b.cfg.BatchSize)

		for {
			select {
			case event := <-b.buffer.ch:
				batch = append(batch, event)
				if len(batch) >= b.cfg.BatchSize {
					b.out <- batch
					batch = make([]MetricEvent, 0, b.cfg.BatchSize)
				}

			case <-ticker.C:
				if len(batch) > 0 {
					b.out <- batch
					batch = make([]MetricEvent, 0, b.cfg.BatchSize)
				}

			case <-b.stop:
				if len(batch) > 0 {
					b.out <- batch
				}
				close(b.out)
				return
			}
		}
	}()
}

func (b *Batcher) Stop() {
	close(b.stop)
}
