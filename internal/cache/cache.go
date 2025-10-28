package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache - простой обёрткой над redis.Client
type Cache struct {
	rdb *redis.Client
}

// New создаёт новый кэш-клиент
func New(addr string) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &Cache{rdb: rdb}
}

// Set сохраняет значение с TTL
func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	return c.rdb.Set(context.Background(), key, value, ttl).Err()
}

// Get возвращает значение по ключу
func (c *Cache) Get(key string) (string, error) {
	return c.rdb.Get(context.Background(), key).Result()
}

// TTL возвращает оставшееся время жизни ключа
func (c *Cache) TTL(key string) (time.Duration, error) {
	return c.rdb.TTL(context.Background(), key).Result()
}
