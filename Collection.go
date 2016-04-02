package mango

import (
	"log"
	"sync"
)

// Collection is map of *Page
type Collection struct {
	sync.RWMutex

	m map[string]PageList
}

// NewCollection - create and init as empty
func NewCollection() *Collection {
	c := &Collection{}
	c.MakeEmpty()
	return c
}

// MakeEmpty - init or clear map
func (c *Collection) MakeEmpty() {
	c.Lock()
	c.m = make(map[string]PageList, 0)
	c.Unlock()
}

// Get from local map by key
func (c *Collection) Get(key string) PageList {
	c.RLock()
	defer c.RUnlock()

	return c.m[key]
}

// Append new page to PageList under key
func (c *Collection) Append(key string, page *Page) {
	c.Lock()
	c.m[key] = append(c.m[key], page)
	c.Unlock()
}

// Remove by key
func (c *Collection) Remove(key string) {
	c.Lock()
	delete(c.m, key)
	c.Unlock()
}

// Print collection items
func (c *Collection) Print(label string) {
	items := c.m
	log.Printf("=== %s (%d) ===============", label, len(items))
	for itemKey, pages := range items {
		log.Printf("- %s: %-12s (%d pages)", label, itemKey, len(pages))
	}
	log.Println()
}
