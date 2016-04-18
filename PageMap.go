package mango

import (
	"log"
	"strings"
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
	pm.RLock()
	defer pm.RUnlock()

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

	log.Printf("--- %d pages ------------------------------------------------", len(pages))
	for slug, p := range pages {
		prefix := ""

		if contentFrom := p.Get("ContentFrom"); contentFrom != "" {
			prefix += contentFrom + " ⏩"
		}

		if p.IsYes("IsUnlisted") {
			prefix += "🔎" // ⛬ ⋱ ⍉ ⛳ 🔎 🔓 🔒 🎩
		}

		if p.IsNo("IsCache") {
			// Not cached
			prefix += " ⟳"
		}

		if redirect := p.Get("Redirect"); redirect != "" {
			if len(redirect) > 20 {
				redirect = strings.TrimPrefix(redirect, "https://")
				redirect = strings.TrimPrefix(redirect, "http://")
				if len(redirect) > 20 {
					redirect = redirect[:18] + ".."
				}
			}
			prefix += redirect + " ⏮"
		}

		if p.IsEqual("Sort", "Reverse") {
			prefix += "[z-a]"
		} else if p.IsEqual("Sort", "Random") {
			prefix += "[?-?]"
		}

		if p.Parent == nil && !p.IsSet("Level") {
			prefix += strings.ToUpper(p.Get("Slug"))
		} else if p.IsEqual("Level", "1") {
			// prefix += "*"
		}

		collectionStr := ""
		for ckey := range p.App.collections {
			if p.IsSet(ckey) {
				collectionStr += "[" + ckey[:1] + "]: "
				cval := p.Get(ckey)
				if len(cval) > 25 {
					// Make shorter and skip middle. Show only start/end values
					v := cval[:5] + " .. " + cval[len(cval)-20:]
					cval = v
				}
				collectionStr += cval + " "
			}
		}

		if p.IsDir() {
			slug = "/" + slug
		}
		log.Printf(" %12s  %-30s %s\n", prefix, slug, collectionStr)
	}
	log.Println()
}
