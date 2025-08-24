package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsServer starts a Prometheus metrics server on the given address.
func StartMetricsServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}

func RunMetricsServer(ctx context.Context, addr string) error {

	errCh := make(chan error, 1)

	go func() {
		errCh <- StartMetricsServer(addr)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	return nil
}

// RegisterCustomMetric registers a custom Prometheus metric.
func RegisterCustomMetric(metric prometheus.Collector) {
	prometheus.MustRegister(metric)
}
