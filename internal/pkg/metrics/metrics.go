package metrics

import (
	"fmt"
	"github.com/YeHeng/go-web-api/internal/pkg/factory"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
)

var metricsRequestsTotal *prometheus.CounterVec
var metricsRequestsCost *prometheus.HistogramVec

func init() {
	factory.Register("metrics", &metricsLifecycle{})
}

type metricsLifecycle struct {
}

func (m *metricsLifecycle) Init() {

	fmt.Println(color.Green("* [prometheus init]"))

	metricsRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "go_web_api",
			Subsystem: "",
			Name:      "requests_total",
			Help:      "request(ms) total",
		},
		[]string{"method", "path"},
	)
	metricsRequestsCost = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "go_web_api",
			Subsystem: "",
			Name:      "requests_cost",
			Help:      "request(ms) cost milliseconds",
		},
		[]string{"method", "path", "success", "http_code", "business_code", "cost_milliseconds", "trace_id"},
	)

	prometheus.MustRegister(metricsRequestsTotal, metricsRequestsCost)
}

func (m *metricsLifecycle) Destroy() {
}

// RecordMetrics 记录指标
func RecordMetrics(method, uri string, success bool, httpCode, businessCode int, costSeconds float64, traceId string) {
	metricsRequestsTotal.With(prometheus.Labels{
		"method": method,
		"path":   uri,
	}).Inc()

	metricsRequestsCost.With(prometheus.Labels{
		"method":            method,
		"path":              uri,
		"success":           cast.ToString(success),
		"http_code":         cast.ToString(httpCode),
		"business_code":     cast.ToString(businessCode),
		"cost_milliseconds": cast.ToString(costSeconds * 1000),
		"trace_id":          traceId,
	}).Observe(costSeconds)
}
