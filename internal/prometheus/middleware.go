package prometheus

import (
	"github.com/best-expendables/router/internal/chi"
	"net/http"
	"strconv"
	"time"

	"github.com/felixge/httpsnoop"

	"github.com/prometheus/client_golang/prometheus"
)

type Counters struct {
	registered bool

	InFlightGauge        prometheus.Gauge
	ReqTotalCounter      *prometheus.CounterVec
	ReqDurationHistogram *prometheus.HistogramVec
}

var counters = &Counters{
	InFlightGauge: prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "http",
		Name:      "in_flight_requests",
		Help:      "A gauge of requests currently being served by the wrapped handler.",
	}),
	ReqTotalCounter: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method", "path"},
	),
	ReqDurationHistogram: prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "A histogram of latencies for requests.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	),
}

// DefaultCounters register counters which uses by default
// Must be called before function: func(h http.Handler) http.Handler
func DefaultCounters() *Counters {
	if !counters.registered {
		prometheus.MustRegister(
			counters.InFlightGauge,
			counters.ReqTotalCounter,
			counters.ReqDurationHistogram,
		)

		counters.registered = true
	}

	return counters
}

func InstrumentHandlerCounter(metric *prometheus.CounterVec, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		metric.With(prometheus.Labels{
			"code":   strconv.Itoa(m.Code),
			"method": r.Method,
			"path":   chi.RoutePatternFromRequest(r),
		}).Inc()
	})
}

func InstrumentHandlerDuration(metric prometheus.ObserverVec, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		metric.With(prometheus.Labels{
			"method": r.Method,
			"path":   chi.RoutePatternFromRequest(r),
		}).Observe(float64(m.Duration / time.Second))
	})
}
