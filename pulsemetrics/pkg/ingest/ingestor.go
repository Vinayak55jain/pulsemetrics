package ingest

import "github.com/vinayak55jain/pulsemetrics/pkg/storage"

type Ingestor struct {
	buffer  *Buffer
	batcher *Batcher
}

func NewIngestor(cfg Config, store storage.Storage) (*Ingestor, error) {
	buffer := NewBuffer(cfg.BufferSize)
	batcher := NewBatcher(buffer, cfg)

	batcher.Start()
	StartWorkers(batcher.out, store, cfg.WorkerCount)

	return &Ingestor{
		buffer:  buffer,
		batcher: batcher,
	}, nil
}

func (i *Ingestor) Push(event MetricEvent) {
	i.buffer.Push(event)
}

func (i *Ingestor) Close() {
	i.batcher.Stop()
}
