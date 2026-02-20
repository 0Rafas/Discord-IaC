package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "bot_messages_received_total",
		Help: "The total number of messages received by the bot",
	})
)

// StartServer starts the prometheus metrics endpoint
func StartServer(addr string) {
	log.Printf("Starting Prometheus metrics server on http://localhost%s/metrics", addr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Metrics server failed: %v", err)
	}
}
