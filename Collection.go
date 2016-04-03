package mango

import (
	"log"
	"strings"
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
	key = c.normalizeKey(key)

	c.RLock()
	defer c.RUnlock()

	return c.m[key]
}

// Append new page to PageList under key
func (c *Collection) Append(key string, page *Page) {
	key = c.normalizeKey(key)

	c.Lock()
	c.m[key] = append(c.m[key], page)
	c.Unlock()
}

// Remove by key
func (c *Collection) Remove(key string) {
	key = c.normalizeKey(key)

	c.Lock()
	delete(c.m, key)
	c.Unlock()
}

// Make key lowercased and trimmed
func (c *Collection) normalizeKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ToLower(key)
	return key
}

// Len is part of sort.Interface.
func (c *Collection) Len() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.m)
}

// Print collection items
func (c *Collection) Print(label string) {
	items := c.m
	log.Printf("--- %s (%d) ---------------------------", label, len(items))
	for itemKey, pages := range items {
		log.Printf("# %-20s (%d pages)", itemKey, len(pages))
	}
}
