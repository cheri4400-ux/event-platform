package metrics

import "github.com/prometheus/client_golang/prometheus"

var EventsProcessed = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "events_processed_total",
		Help: "Total processed events",
	},
)

var EventsFailed = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "events_failed_total",
		Help: "Total failed events",
	},
)

func Init() {
	prometheus.MustRegister(EventsProcessed)
	prometheus.MustRegister(EventsFailed)
}
