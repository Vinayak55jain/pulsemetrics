package ingest

import model "github.com/vinayak55jain/pulsemetrics/pkg/model"

// MetricEvent is an alias to the shared model type. Using a shared model
// package removes import cycles between ingest and storage while keeping the
// ingest package API stable and simple to use.
type MetricEvent = model.MetricEvent
