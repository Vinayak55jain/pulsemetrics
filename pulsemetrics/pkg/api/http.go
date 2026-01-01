package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vinayak55jain/pulsemetrics/pkg/ingest"
)

func RegisterRoutes(mux *http.ServeMux, ing *ingest.Ingestor) {
	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		var event ingest.MetricEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}
		ing.Push(event) // NON-BLOCKING

		w.WriteHeader(http.StatusAccepted)
	})

	// Simple JSON metrics endpoint exposing the ingestion counters.
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ms := ingest.SnapshotMetrics()
		_ = json.NewEncoder(w).Encode(ms)
	})
}
