package ingest

import (
    "context"
    "sync"
    "testing"
    "time"

    m "github.com/vinayak55jain/pulsemetrics/pkg/model"
    "github.com/vinayak55jain/pulsemetrics/pkg/storage"
)

type testStore struct {
    mu      sync.Mutex
    batches [][]m.MetricEvent
}

func (s *testStore) Save(batch []m.MetricEvent) error {
    time.Sleep(20 * time.Millisecond)
    s.mu.Lock()
    defer s.mu.Unlock()
    copied := make([]m.MetricEvent, len(batch))
    copy(copied, batch)
    s.batches = append(s.batches, copied)
    return nil
}

func (s *testStore) Count() int {
    s.mu.Lock()
    defer s.mu.Unlock()
    c := 0
    for _, b := range s.batches {
        c += len(b)
    }
    return c
}

func TestIngestorBufferingAndBatching(t *testing.T) {
    cfg := DefaultConfig()
    cfg.BatchSize = 3
    cfg.BufferSize = 5
    cfg.WorkerCount = 1
    cfg.FlushInterval = 50 * time.Millisecond

    store := &testStore{}

    ing, err := NewIngestor(cfg, storage.Storage(store))
    if err != nil {
        t.Fatalf("NewIngestor error: %v", err)
    }
    defer ing.Close()

    // Push 7 events quickly with buffer size 5 some events may be dropped.
    totalPush := 7
    for i := 0; i < totalPush; i++ {
        ev := m.MetricEvent{Type: "test", Name: "m", Value: float64(i), Timestamp: time.Now()}
        ing.Push(MetricEvent(ev))
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    ticker := time.NewTicker(20 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            t.Fatalf("timeout waiting for store to receive batches; received=%d", store.Count())
        case <-ticker.C:
            if store.Count() >= 1 {
                time.Sleep(100 * time.Millisecond)
                t.Logf("store received %d metrics (of %d pushed)", store.Count(), totalPush)
                return
            }
        }
    }
}
