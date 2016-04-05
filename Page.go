package mango

import (
	"log"
	"strconv"
	"strings"
	"sync"
)

// Page - page with content and params + sub-pages
type Page struct {
	sync.RWMutex

	// Linkt to application
	App *Application

	// Content
	content []byte

	// Params that describe this page
	Params map[string]string

	// Parent page
	Parent *Page

	// Sub-pages for this page
	Pages PageList
}

// newPage - create page from label
// Other params can be set after return *Page
// client must use app.NewPage("Label") to create new virtual page
func newPage(label string) *Page {
	page := &Page{
		// We creating pseuode Page (not exists on filesystem)
		// so need to make it look like filename so it can be parsed properly
		Params: filenameToParams(label + _Md),
	}

	// Mark that this page is create not from file
	page.Params["IsVirtual"] = _Yes

	// Slug is used for real pages
	page.Params["VirtualSlug"] = page.Params["Slug"]
	delete(page.Params, "Slug")

	return page
}

// fileToPage - create/init new page from existing file
func fileToPage(fpath string) *Page {

	// Extract content
	params := fileToParams(fpath)
	bufContent := []byte(params["Content"])
	delete(params, "Content")

	// Create new page
	page := newPage("")
	page.SetContent(bufContent)
	page.Params = params // assign original params

	return page
}

// SetContent set content for page
func (page *Page) SetContent(content []byte) {
	page.Lock()
	page.content = content
	page.Unlock()
}

// AppendContent - append to content
func (page *Page) AppendContent(content []byte) {
	pageContent := page.Content()
	page.SetContent(append(pageContent, content...))
}

// Content - get content for page
func (page *Page) Content() []byte {
	page.RLock()
	defer page.RUnlock()

	return page.content
}

// Set - set thread-safely param to Page.Params
func (page *Page) Set(key, val string) {
	if key == "Slug" {
		// Slug can't be changed after loading all pages
		// If slug must be changed:
		// - rename file.md
		// - add Slug: param to file.md header section
		return
	}

	page.Lock()
	page.Params[key] = val
	page.Unlock()

}

// Get - get thread-safely param to Page.Params
func (page *Page) Get(key string) string {
	page.RLock()
	defer page.RUnlock()

	return page.Params[key]
}

// Split - get param as slice splitted by given separator
func (page *Page) Split(key, sep string) []string {
	val := page.Get(key)

	// A, B,,,,C --> results in 3 items
	arr := strings.Split(val, sep) // dirty list
	var _arr []string              // validated list
	for _, v := range arr {
		v := strings.TrimSpace(v)
		if v != "" {
			// Only with content are added
			_arr = append(_arr, v)
		}
	}

	return _arr
}

// IsEqual - shorthand to compare param with custom string
func (page *Page) IsEqual(key, val string) bool {
	return page.Get(key) == val
}

// IsYes - shorthand to compare param with "Yes"
func (page *Page) IsYes(key string) bool {
	return page.IsEqual(key, _Yes)
}

// IsNo - shorthand to compare param with "No"
func (page *Page) IsNo(key string) bool {
	return page.IsEqual(key, _No)
}

// IsSet - shorthand to find out is this val set and not empty ""
func (page *Page) IsSet(key string) bool {
	return !page.IsEqual(key, "")
}

// IsDir - shorthand to find out is this val set and not empty "IsDir"
func (page *Page) IsDir() bool {
	return page.IsYes("IsDir")
}

// Check if page is duplicate slug
func (page *Page) isDuplicate() bool {
	return page.App.Page(page.Get("Slug")) != nil
}

// Get some params from path
func (page *Page) setPathParams() {

	// relative path from app.ContentPath
	rpath := strings.TrimPrefix(page.Get("Path"), page.App.ContentPath)

	// Remove filename
	rpath = strings.TrimSuffix(rpath, page.Get("FileName"))

	// split to parts
	rpath = strings.Trim(rpath, "/")
	arr := strings.Split(rpath, "/")

	// level of depth
	if len(arr) == 1 && arr[0] == "" {
		// langage is in zero level depth
		// remove empty
		arr = make([]string, 0)
	}

	// Need at least 1
	if len(arr) == 0 {
		return
	}

	// Set Level of depth
	page.Set("Lang", arr[0])
	page.Set("Level", strconv.Itoa(len(arr)))

	// 1. en -> 2. top-menu -> 3-n.pages...
	// 2. is group keys. Every language folder have same groupkeys
	// so we need to prefix these slugs with language
	// en-top-menu
	if page.IsEqual("Level", "1") && page.IsDir() {
		page.Params["Slug"] = page.Get("Lang") + "-" + page.Get("Slug")
	}

	// Need at least 2
	if len(arr) < 2 {
		return
	}

	// Set params based on arr
	page.Set("GroupKey", arr[1])
}

// Generate unique slug based on old one
func (page *Page) avoidDuplicate() {
	// Suffix loop by count until unique
	origSlug := page.Get("Slug")
	for i := 2; page.isDuplicate(); i++ {
		page.Lock()
		page.Params["Slug"] = origSlug + "-" + strconv.Itoa(i)
		page.Unlock()
	}

}

// Walk all down by sub-pages and do custom stuff
// Can be customized by custom func
// TODO: goroutines?
func (page *Page) Walk(fnCheck func(p *Page) bool) PageList {
	pages := make(PageList, 0)

	page.RLock()
	all := page.Pages[:]
	page.RUnlock()

	for _, p := range all {
		if fnCheck(p) {
			pages = append(pages, p)
		}

		// Go deeper
		if p.IsDir() {
			pages = append(pages, p.Walk(fnCheck)...)
		}
	}

	return pages
}

// Search - find all pages by given search term
// TODO: make correct search by params and content
func (page *Page) Search(sterm string) PageList {
	sterm = strings.TrimSpace(sterm)
	sterm = strings.ToLower(sterm)

	return page.Walk(func(p *Page) bool {
		// TODO: skip unlisted pages?

		// Custom check
		// TODO: add correct search by params and content. Not only slug
		s := p.Get("Slug") +
			p.Get("Label") +
			p.Get("Title") +
			string(p.Content())
		s = strings.ToLower(s)

		isFound := strings.Index(s, sterm) >= 0
		return isFound
	})
}

// SearchByParam - find all pages that search value is equal to page param
func (page *Page) SearchByParam(key, val string) PageList {
	return page.Walk(func(p *Page) bool {
		// Check for equal param values
		return p.IsEqual(key, val)
	})
}

// Print pages in list
func (page *Page) Print() {
	printMap(page.Get("Slug"), page.Params)
}

// PrintTree - Print all pages under this page
func (page *Page) PrintTree(depth int) {
	for _, p := range page.Pages {
		log.Printf("%s %-30s %-30s %3d bytes", strings.Repeat("    ", depth), p.Params["Label"], p.Params["Slug"], len(p.Content()))

		// printMap(p.Params["Label"], p.Params)
		if len(p.Pages) > 0 {
			p.PrintTree(depth + 1)
		}
	}
}

// MergeParams - merge some more params
func (page *Page) MergeParams(moreParams map[string]string) {
	page.RLock()
	pageParams := page.Params
	page.RUnlock()

	pageParams = mergeParams(pageParams, moreParams)

	page.Lock()
	page.Params = pageParams
	page.Unlock()
}
