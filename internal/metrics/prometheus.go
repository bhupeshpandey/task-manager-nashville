package metrics

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-nashville/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	registry *prometheus.Registry
	config   *models.MetricsConfig
	logger   models.Logger
)

func InitializePrometheus(metricsConfig *models.MetricsConfig) {
	config = metricsConfig
	logger = metricsConfig.Logger
	registry = prometheus.NewRegistry()
}

// RegisterMetric this func allows other sub services to register their metrics with prometheus
func RegisterMetric(metric prometheus.Collector) {
	registry.MustRegister(metric)
}

func StartServer() {
	// Expose the metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	logger.Log(models.InfoLevel, "Starting the prometheus server", fmt.Sprintf("%s:%s", config.Host, config.Port))
	// Start the HTTP server
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
	if err != nil {
		logger.Log(models.ErrorLevel, "Unable to start prometheus server", err)
		return
	}
}
