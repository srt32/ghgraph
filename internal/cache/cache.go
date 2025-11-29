package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// entry represents a cached value with expiration
type entry struct {
	value      interface{}
	expiration time.Time
}

// Cache is a thread-safe in-memory cache with TTL and token scoping
type Cache struct {
	data sync.Map
	ttl  time.Duration
}

// NewCache creates a new cache with the specified TTL
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		ttl: ttl,
	}
	// Start cleanup goroutine
	go c.cleanup()
	return c
}

// makeKey creates a cache key scoped by token
func (c *Cache) makeKey(token, key string) string {
	// Hash token for privacy and consistent key length
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:8]) // Use first 8 bytes
	return tokenHash + ":" + key
}

// Get retrieves a value from the cache
// Returns nil if not found or expired
func (c *Cache) Get(token, key string) interface{} {
	fullKey := c.makeKey(token, key)
	val, ok := c.data.Load(fullKey)
	if !ok {
		return nil
	}

	e := val.(*entry)
	if time.Now().After(e.expiration) {
		c.data.Delete(fullKey)
		return nil
	}

	return e.value
}

// Set stores a value in the cache with TTL
func (c *Cache) Set(token, key string, value interface{}) {
	fullKey := c.makeKey(token, key)
	c.data.Store(fullKey, &entry{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	})
}

// Delete removes a value from the cache
func (c *Cache) Delete(token, key string) {
	fullKey := c.makeKey(token, key)
	c.data.Delete(fullKey)
}

// Clear removes all entries for a specific token
func (c *Cache) Clear(token string) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:8])
	prefix := tokenHash + ":"

	c.data.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				c.data.Delete(key)
			}
		}
		return true
	})
}

// cleanup periodically removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.data.Range(func(key, value interface{}) bool {
			e := value.(*entry)
			if now.After(e.expiration) {
				c.data.Delete(key)
			}
			return true
		})
	}
}
