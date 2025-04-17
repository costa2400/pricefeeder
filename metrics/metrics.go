package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const PrometheusNamespace = "pricefeeder"

// Basic Operational Metrics

// PriceSourceCounter tracks the total number of price fetches by source and success status
var PriceSourceCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: PrometheusNamespace,
	Name:      "fetched_prices_total",
	Help:      "The total number prices fetched, by source and success status",
}, []string{"source", "success"})

// AggregatePriceCounter tracks the number of price aggregations by pair, source, and success status
var AggregatePriceCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: PrometheusNamespace,
	Name:      "aggregate_prices_total",
	Help:      "The total number of times prices were aggregated by pair, source, and success status",
}, []string{"pair", "source", "success"})

// PostedPricesCounter tracks the number of posted prices by success status
var PostedPricesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: PrometheusNamespace,
	Name:      "prices_posted_total",
	Help:      "The total number of txs sent to the on-chain oracle module",
}, []string{"success"})

// PriceFetchLatency tracks how long it takes to fetch prices from each source
var PriceFetchLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: PrometheusNamespace,
	Name:      "price_fetch_latency_seconds",
	Help:      "The time it takes to fetch prices from each source in seconds",
	Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0},
}, []string{"source", "pair"})

// TxBroadcastLatency tracks how long it takes to broadcast transactions
var TxBroadcastLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: PrometheusNamespace,
	Name:      "tx_broadcast_latency_seconds",
	Help:      "The time it takes to broadcast transactions in seconds",
	Buckets:   []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 30.0},
}, []string{"tx_type"})

// Data Quality Metrics

// PriceDeviation tracks the deviation between consecutive price updates
var PriceDeviation = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "price_deviation_percent",
	Help:      "The percentage deviation between consecutive price updates",
}, []string{"pair", "source"})

// DataFreshness tracks how old the latest price data is
var DataFreshness = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "data_freshness_seconds",
	Help:      "How old the latest price data is in seconds",
}, []string{"pair", "source"})

// CrossSourceDeviation tracks the deviation in prices between different sources for the same pair
var CrossSourceDeviation = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "cross_source_deviation_percent",
	Help:      "The percentage deviation in prices between different sources for the same pair",
}, []string{"pair", "source_primary", "source_secondary"})

// System Health Metrics

// ErrorCount tracks the number of errors by type and component
var ErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: PrometheusNamespace,
	Name:      "error_count_total",
	Help:      "The total number of errors by type and component",
}, []string{"error_type", "component"})

// GoroutineCount tracks the number of active goroutines
var GoroutineCount = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "goroutine_count",
	Help:      "The number of active goroutines",
})

// MemoryUsage tracks the memory usage of the application
var MemoryUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "memory_usage_bytes",
	Help:      "The memory usage of the application in bytes",
}, []string{"type"}) // type can be "alloc", "sys", etc.

// ConnectionStatus tracks the status of connections to external services
var ConnectionStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "connection_status",
	Help:      "The status of connections to external services (1 for connected, 0 for disconnected)",
}, []string{"service_type", "endpoint"})

// RateLimitStatus tracks the rate limit status for external APIs
var RateLimitStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: PrometheusNamespace,
	Name:      "rate_limit_remaining",
	Help:      "The remaining rate limit for external APIs",
}, []string{"service", "endpoint"})
