package cache

import (
	"awesomeProject/internal/config"
	"sync"
	"time"
)

// CacheItem представляет элемент кэша с временем истечения
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache реализует кэш с TTL
type Cache struct {
	items    sync.Map
	config   *config.Config
	stopChan chan struct{}
}

var (
	instance *Cache
	once     sync.Once
)

// GetCache возвращает синглтон кэша
func GetCache() *Cache {
	once.Do(func() {
		instance = &Cache{
			config:   config.GetConfig(),
			stopChan: make(chan struct{}),
		}
		if instance.config.Cache.Enabled {
			go instance.cleanupLoop()
		}
	})
	return instance
}

// Set сохраняет значение в кэш
func (c *Cache) Set(key string, value interface{}) {
	if !c.config.Cache.Enabled {
		return
	}

	item := CacheItem{
		Value:      value,
		Expiration: time.Now().Add(c.config.Cache.TTL),
	}
	c.items.Store(key, item)
}

// Get получает значение из кэша
func (c *Cache) Get(key string) (interface{}, bool) {
	if !c.config.Cache.Enabled {
		return nil, false
	}

	if value, ok := c.items.Load(key); ok {
		item := value.(CacheItem)
		if time.Now().Before(item.Expiration) {
			return item.Value, true
		}
		// Удаляем просроченный элемент
		c.items.Delete(key)
	}
	return nil, false
}

// Delete удаляет значение из кэша
func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

// Clear очищает весь кэш
func (c *Cache) Clear() {
	c.items = sync.Map{}
}

// cleanupLoop периодически очищает просроченные элементы
func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.config.Cache.CleanupTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopChan:
			return
		}
	}
}

// cleanup удаляет все просроченные элементы
func (c *Cache) cleanup() {
	now := time.Now()
	c.items.Range(func(key, value interface{}) bool {
		item := value.(CacheItem)
		if now.After(item.Expiration) {
			c.items.Delete(key)
		}
		return true
	})
}

// Stop останавливает очистку кэша
func (c *Cache) Stop() {
	close(c.stopChan)
}

// GetSize возвращает текущий размер кэша
func (c *Cache) GetSize() int {
	var size int
	c.items.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}
