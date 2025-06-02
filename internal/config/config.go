package config

import (
	"os"
	"strconv"
	"time"
)

// Config содержит все настройки приложения
type Config struct {
	// Настройки таймаутов
	Timeout struct {
		NewsDetails    time.Duration
		UserDetails    time.Duration
		DatabaseQuery  time.Duration
		CacheDuration  time.Duration
		RequestTimeout time.Duration
	}

	// Настройки кэширования
	Cache struct {
		Enabled     bool
		MaxSize     int
		TTL         time.Duration
		CleanupTime time.Duration
	}

	// Настройки метрик
	Metrics struct {
		Enabled     bool
		Port        int
		Path        string
		PushGateway string
	}
}

var cfg *Config

// GetConfig возвращает текущую конфигурацию
func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
		cfg.loadFromEnv()
	}
	return cfg
}

// loadFromEnv загружает настройки из переменных окружения
func (c *Config) loadFromEnv() {
	// Таймауты
	c.Timeout.NewsDetails = getDurationEnv("NEWS_DETAILS_TIMEOUT", 500*time.Millisecond)
	c.Timeout.UserDetails = getDurationEnv("USER_DETAILS_TIMEOUT", 500*time.Millisecond)
	c.Timeout.DatabaseQuery = getDurationEnv("DB_QUERY_TIMEOUT", 1*time.Second)
	c.Timeout.CacheDuration = getDurationEnv("CACHE_DURATION", 5*time.Minute)
	c.Timeout.RequestTimeout = getDurationEnv("REQUEST_TIMEOUT", 5*time.Second)

	// Кэширование
	c.Cache.Enabled = getBoolEnv("CACHE_ENABLED", true)
	c.Cache.MaxSize = getIntEnv("CACHE_MAX_SIZE", 1000)
	c.Cache.TTL = getDurationEnv("CACHE_TTL", 5*time.Minute)
	c.Cache.CleanupTime = getDurationEnv("CACHE_CLEANUP_TIME", 10*time.Minute)

	// Метрики
	c.Metrics.Enabled = getBoolEnv("METRICS_ENABLED", true)
	c.Metrics.Port = getIntEnv("METRICS_PORT", 9090)
	c.Metrics.Path = getStringEnv("METRICS_PATH", "/metrics")
	c.Metrics.PushGateway = getStringEnv("METRICS_PUSH_GATEWAY", "")
}

// Вспомогательные функции для получения значений из переменных окружения
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getStringEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
