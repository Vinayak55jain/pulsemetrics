package ingest

import (
	"context"
	"sync"

	m "github.com/vinayak55jain/pulsemetrics/pkg/model"
	"github.com/vinayak55jain/pulsemetrics/pkg/storage"
)

type Ingestor struct {
	cfg     Config
	buffer  *Buffer
	batcher *Batcher

	// background control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// DLQ
	dlqCh   chan []m.MetricEvent
	dlqFile string
}

func NewIngestor(cfg Config, store storage.Storage) (*Ingestor, error) {
	buffer := NewBuffer(cfg.BufferSize)
	batcher := NewBatcher(buffer, cfg)

	ctx, cancel := context.WithCancel(context.Background())

	ing := &Ingestor{
		cfg:     cfg,
		buffer:  buffer,
		batcher: batcher,
		ctx:     ctx,
		cancel:  cancel,
		dlqCh:   make(chan []m.MetricEvent, 100),
		dlqFile: "dlq.jsonl",
	}

	ing.batcher.Start()

	StartWorkers(ing.ctx, &ing.wg, ing.batcher.out, store, cfg.WorkerCount, ing.dlqCh, ing.dlqFile)

	ing.wg.Add(1)
	go func() {
		defer ing.wg.Done()
		for {
			select {
			case <-ing.ctx.Done():
				return
			case b, ok := <-ing.dlqCh:
				if !ok {
					return
				}
				// simply log DLQ batches; they are also persisted to file by workers
				_ = b // placeholder; production code could push to monitoring
			}
		}
	}()

	return ing, nil
}

func (i *Ingestor) Push(event MetricEvent) {
	i.buffer.Push(event)
}

func (i *Ingestor) Close() {
	i.cancel()
	i.batcher.Stop()
	i.wg.Wait()
	i.buffer.Close()
	close(i.dlqCh)
}
