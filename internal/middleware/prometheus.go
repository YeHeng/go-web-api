package middleware

import (
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cast"
)

var metricsRequestsTotal *prometheus.CounterVec
var metricsRequestsCost *prometheus.HistogramVec

var prometheusHandler gin.HandlerFunc

func init() {
	AddMiddleware(&prometheusMiddleware{})
}

type prometheusMiddleware struct {
}

func (m *prometheusMiddleware) Init() {

	cfg := config.Get().Feature
	if cfg.RecordMetrics {

		fmt.Println(color.Green("* [register middleware metrics]"))

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

		prometheusHandler = gin.WrapH(promhttp.Handler())

	}

}

func (m *prometheusMiddleware) Apply(r *gin.Engine) {
	cfg := config.Get().Feature
	if cfg.RecordMetrics {
		r.GET("/metrics", prometheusHandler) // register prometheus
	}
}

func (m *prometheusMiddleware) Get() gin.HandlerFunc {
	return prometheusHandler
}

func (m *prometheusMiddleware) Destroy() {
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
