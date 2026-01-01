package ingest

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	m "github.com/vinayak55jain/pulsemetrics/pkg/model"
	"github.com/vinayak55jain/pulsemetrics/pkg/storage"
)

// StartWorkers starts a bounded pool of workers that process batches
// asynchronously. It accepts a context for cancellation, a WaitGroup to
// track worker lifecycle, a dead-letter channel for failed batches and a
// path to a DLQ file where failed batches will be appended as JSON lines.
func StartWorkers(ctx context.Context, wg *sync.WaitGroup, batches <-chan []MetricEvent, store storage.Storage, workers int, dlqCh chan<- []m.MetricEvent, dlqFile string) {
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case batch, ok := <-batches:
					if !ok {
						return
					}
					// Attempt to store with retries
					var err error
					for attempt := 0; attempt < 3; attempt++ {
						err = store.Save(batch)
						if err == nil {
							break
						}
						// small backoff
						time.Sleep(10 * time.Millisecond)
					}
					if err != nil {
						log.Printf("worker %d: failed to store batch after retries: %v", id, err)
						// send to DLQ channel if available
						if dlqCh != nil {
							// convert to model.MetricEvent for DLQ contract
							mBatch := make([]m.MetricEvent, len(batch))
							for i := range batch {
								mBatch[i] = m.MetricEvent(batch[i])
							}
							select {
							case dlqCh <- mBatch:
							default:
								log.Println("dlq channel full, dropping dlq batch")
							}
						}
						// append to DLQ file for inspection
						if dlqFile != "" {
							f, ferr := os.OpenFile(dlqFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
							if ferr != nil {
								log.Printf("failed to open dlq file: %v", ferr)
								continue
							}
							enc := json.NewEncoder(f)
							if err := enc.Encode(batch); err != nil {
								log.Printf("failed to write dlq batch to file: %v", err)
							}
							f.Close()
						}
					}
				}
			}
		}(i)
	}
}
