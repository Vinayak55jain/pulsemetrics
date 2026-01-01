package ingest

import "github.com/vinayak55jain/pulsemetrics/pkg/storage"

func StartWorkers(
	batches <-chan []MetricEvent,
	store storage.Storage,
	workers int,
) {
	for i := 0; i < workers; i++ {
		go func(id int) {
			for batch := range batches {
				store.Save(batch)
			}
		}(i)
	}
}
