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

// NewPageMap - create and init as empty
func NewPageMap() *PageMap {
	pm := &PageMap{}
	pm.MakeEmpty()
	return pm
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
// Can't use only Slug, because PageMap can be used for many purposes
func (pm *PageMap) Add(key string, page *Page) {
	// key := page.Get("Slug")
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

// Filter - filter by custom func
func (pm *PageMap) Filter(fnCheck func(p *Page) bool) PageList {
	pages := make(PageList, 0)

	pm.RLock()
	for _, p := range pm.m {
		if p != nil && fnCheck(p) {
			pages = append(pages, p)
		}
	}
	pm.RUnlock()

	return pages
}

// Print pages in list
func (pm *PageMap) Print() {
	pages := pm.m

	if len(pages) == 0 {
		return
	}

	log.Println("------------------------------------------------------------")
	for slug, p := range pages {
		dir := ""
		redirect := ""
		contentFrom := ""

		if p.IsDir() {
			dir = "/"
		}

		// if p.IsSet("Redirect") {
		// 	redirect = "R"
		// }
		//
		// if p.IsSet("ContentFrom") {
		// 	contentFrom = "^"
		// }

		log.Printf("- %-24s %2s%2s%2s%2s %24s\n", slug, dir, redirect, contentFrom, p.Get("Level"), p.Get("Title"))
	}
	log.Println("------------------------------------------------------------")
}
