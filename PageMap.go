package mango

import (
	"log"
	"sync"
)

// PageMap is map of *Page
type PageMap struct {
	sync.RWMutex

	m map[string]*Page
}

// MakeEmpty - init or clear map
func (pm *PageMap) MakeEmpty() {
	pm.Lock()
	pm.m = make(map[string]*Page, 0)
	pm.Unlock()
}

// Get from local map by key
func (pm *PageMap) Get(key string) *Page {
	pm.RLock()
	defer pm.RUnlock()

	return pm.m[key]
}

// Add new *Page to local map
func (pm *PageMap) Add(page *Page) {
	key := page.Get("Slug")

	pm.Lock()
	pm.m[key] = page
	pm.Unlock()
}

// Remove by key
func (pm *PageMap) Remove(key string) {
	pm.Lock()
	delete(pm.m, key)
	pm.Unlock()
}

// Len - map item count
func (pm *PageMap) Len() int {
	return len(pm.m)
}

// Print pages in list
func (pm *PageMap) Print() {
	pages := pm.m

	log.Println("---------------------------------------------------")
	for slug, p := range pages {
		log.Printf("- %20s - %s \n", slug, p.Get("Title"))
	}
	log.Println("---------------------------------------------------")
}
