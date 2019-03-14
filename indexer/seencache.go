package indexer

import (
	"fmt"
	"crypto/sha1"
	"time"
)

type SeenCache struct {
	seen map[string]int
}

func NewSeenCache() *SeenCache {
	return &SeenCache{seen: make(map[string]int)}
}

func (c *SeenCache) Hash(filePath string, fileSize int64, modTime time.Time) string {
	token := fmt.Sprintf("%08d|%s|%s",fileSize,modTime,filePath)
	return fmt.Sprintf("%x", sha1.Sum([]byte(token)))[:10]
}

func (c *SeenCache) Clear() {
	c.seen = make(map[string]int)
}

func (c *SeenCache) Add(filePath string, fileSize int64, modTime time.Time, index int) string {
	hash := c.Hash(filePath, fileSize, modTime)
	c.seen[hash] = index
	return hash
}

func (c *SeenCache) AddHash(hash string, index int) {
	c.seen[hash] = index
}

func (c *SeenCache) Seen(filePath string, fileSize int64, modTime time.Time) (bool, string, int) {
	hash := c.Hash(filePath, fileSize, modTime)
	if index, exists := c.seen[hash]; exists {
		return true, hash, index
	} else {
		return false, "", -1
	}
}

func (c *SeenCache) SeenHash(hash string) (bool, int) {
	index, exists := c.seen[hash]
	return exists, index
}