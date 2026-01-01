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

		event.Timestamp = time.Now()
		ing.Push(event) // NON-BLOCKING

		w.WriteHeader(http.StatusAccepted)
	})
}
