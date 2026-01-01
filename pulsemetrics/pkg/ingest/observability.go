package ingest

import "sync/atomic"

// Basic observability counters. We keep them as package-level atomics for
// simplicity. They are updated by components (buffer, batcher, workers) and
// exposed via the /metrics endpoint.
var (
	metricsDropped uint64 // number of metrics dropped due to full buffer
	batchesFlushed uint64 // number of batches emitted by the batcher
	ingestionCount uint64 // number of metrics accepted into buffer
)

// snapshot returns current metric values. Useful for HTTP /metrics handler.
type MetricsSnapshot struct {
	MetricsDropped uint64 `json:"metrics_dropped"`
	BatchesFlushed uint64 `json:"batches_flushed"`
	Ingested       uint64 `json:"ingested"`
}

func snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		MetricsDropped: atomic.LoadUint64(&metricsDropped),
		BatchesFlushed: atomic.LoadUint64(&batchesFlushed),
		Ingested:       atomic.LoadUint64(&ingestionCount),
	}
}

// SnapshotMetrics exposes the current metrics snapshot for external callers.
func SnapshotMetrics() MetricsSnapshot { return snapshot() }
