package metrics

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/config"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Метрики для новостей
	newsRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "news_request_duration_seconds",
			Help:    "Duration of news requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	newsRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "news_requests_total",
			Help: "Total number of news requests",
		},
		[]string{"endpoint", "method", "status"},
	)

	// Метрики для пользователей
	userRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_request_duration_seconds",
			Help:    "Duration of user requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	userRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_requests_total",
			Help: "Total number of user requests",
		},
		[]string{"endpoint", "method", "status"},
	)

	// Метрики кэша
	cacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	cacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Current size of the cache",
		},
	)
)

// MetricsMiddleware для Gin
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Продолжаем обработку запроса
		c.Next()

		// Записываем метрики
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// Определяем тип запроса и записываем соответствующие метрики
		switch {
		case path == "/protected/news/all" || path == "/protected/news/all-with-details":
			newsRequestDuration.WithLabelValues(path, method).Observe(duration)
			newsRequestTotal.WithLabelValues(path, method, string(status)).Inc()
		case path == "/protected/user/all" || path == "/protected/user/name/:id":
			userRequestDuration.WithLabelValues(path, method).Observe(duration)
			userRequestTotal.WithLabelValues(path, method, string(status)).Inc()
		}
	}
}

// StartMetricsServer запускает сервер метрик
func StartMetricsServer() {
	cfg := config.GetConfig()
	if !cfg.Metrics.Enabled {
		return
	}

	http.Handle(cfg.Metrics.Path, promhttp.Handler())
	go func() {
		addr := ":" + strconv.Itoa(cfg.Metrics.Port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
}

// RecordCacheHit записывает попадание в кэш
func RecordCacheHit() {
	cacheHits.Inc()
}

// RecordCacheMiss записывает промах кэша
func RecordCacheMiss() {
	cacheMisses.Inc()
}

// UpdateCacheSize обновляет размер кэша
func UpdateCacheSize(size int) {
	cacheSize.Set(float64(size))
}

// MetricsCollector для сбора метрик из кэша
type MetricsCollector struct {
	cache *cache.Cache
}

// NewMetricsCollector создает новый коллектор метрик
func NewMetricsCollector(cache *cache.Cache) *MetricsCollector {
	return &MetricsCollector{cache: cache}
}

// Collect собирает метрики из кэша
func (mc *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	// Используем публичный метод для получения размера кэша
	// Предполагаем, что у Cache есть метод GetSize()
	size := mc.cache.GetSize()
	cacheSize.Set(float64(size))
}

// Describe описывает метрики
func (mc *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(mc, ch)
}
