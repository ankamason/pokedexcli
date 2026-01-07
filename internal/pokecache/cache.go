package pokecache

import (
    "sync"
    "time"
)

// cacheEntry represents a single cached item
type cacheEntry struct {
    createdAt time.Time
    val       []byte
}

// Cache is a thread-safe cache with automatic expiration
type Cache struct {
    entries  map[string]cacheEntry
    mutex    sync.Mutex
    interval time.Duration
}

// NewCache creates a new cache with the given expiration interval
func NewCache(interval time.Duration) *Cache {
    c := &Cache{
        entries:  make(map[string]cacheEntry),
        interval: interval,
    }

    // Start the reap loop in a goroutine
    go c.reapLoop()

    return c
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    c.entries[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
}

// Get retrieves an entry from the cache
// Returns the value and a bool indicating if it was found
func (c *Cache) Get(key string) ([]byte, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    entry, exists := c.entries[key]
    if !exists {
        return nil, false
    }

    return entry.val, true
}

// reapLoop runs continuously and removes expired entries
func (c *Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()

    for range ticker.C {
        c.reap()
    }
}

// reap removes entries older than the interval
func (c *Cache) reap() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    now := time.Now()
    for key, entry := range c.entries {
        if now.Sub(entry.createdAt) > c.interval {
            delete(c.entries, key)
        }
    }
}
