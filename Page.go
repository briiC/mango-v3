package mango

import (
	"fmt"
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
	Content []byte

	// Params that describe this page
	Params map[string]string

	// Parent page
	Parent *Page

	// Sub-pages for this page
	Pages PageList
}

// newPage - create/init new page
func newPage(app *Application, fpath string) *Page {
	// TODO: add page to app-wide page list with its slug?
	// so it could be checked (putside this func) for existance?

	// Extract content
	params := fileToParams(fpath)
	bufContent := []byte(params["Content"])
	delete(params, "Content")

	// Create new page
	page := &Page{
		App:     app,
		Content: bufContent,
		Params:  params,
	}

	// some params from path
	page.pathToParams()

	return page
}

// Set - set thread-safely param to Page.Params
func (page *Page) Set(key, val string) {
	if key == "Slug" {
		// TODO: on Slug change need to change this also in app.pageList
		// Because pageList[slug] using slug as index/key
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

// IsEqual - shorthand to compare param with custom string
func (page *Page) IsEqual(key, val string) bool {
	return page.Get(key) == val
}

// IsYes - shorthand to compare param with "Yes"
func (page *Page) IsYes(key string) bool {
	return page.IsEqual(key, "Yes")
}

// IsNo - shorthand to compare param with "No"
func (page *Page) IsNo(key string) bool {
	return page.IsEqual(key, "No")
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
	_, isDuplicate := page.App.pageList[page.Params["Slug"]]
	return isDuplicate
}

// Get some params from path
func (page *Page) pathToParams() {

	// relative path from app.ContentPath
	rpath := strings.TrimPrefix(page.Params["Path"], page.App.ContentPath)

	// Remove filename
	rpath = strings.TrimSuffix(rpath, page.Params["FileName"])

	// split to parts
	rpath = strings.Trim(rpath, "/")
	arr := strings.Split(rpath, "/")

	// level of depth
	if len(arr) == 1 && arr[0] == "" {
		// langage is in zero level depth
		// remove empty
		arr = make([]string, 0)
	}

	// Set Level of depth
	page.Params["Level"] = strconv.Itoa(len(arr))

	if len(arr) < 2 {
		return
	}

	// Set params based on arr
	page.Params["Lang"] = arr[0]
	page.Params["GroupKey"] = arr[1]
}

// Generate unique slug based on old one
func (page *Page) avoidDuplicate() {
	// Suffix loop by count until unique
	origSlug := page.Params["Slug"]
	for i := 2; page.isDuplicate(); i++ {
		page.Params["Slug"] = origSlug + "-" + strconv.Itoa(i)
	}

}

// Walk all down by sub-pages and do custom stuff
// Can be customized by custom func
func (page *Page) Walk(fnCheck func(p *Page) bool) PageList {
	pages := make(PageList, 0)

	for _, p := range page.Pages {
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
func (page *Page) Search(s string) PageList {
	return page.Walk(func(p *Page) bool {
		// Custom check
		// TODO: add correct search by params and content. Not only slug
		return strings.Index(p.Params["Slug"], s) >= 0
	})
}

// SearchByParam - find all pages that search value is equal to page param
func (page *Page) SearchByParam(key, val string) PageList {
	return page.Walk(func(p *Page) bool {
		// Check for equal param values
		return p.IsEqual(key, val)
	})
}

// PrintTree - Print all pages under this page
func (page *Page) PrintTree(depth int) {
	for _, p := range page.Pages {
		fmt.Printf("%s %-30s %-30s", strings.Repeat("    ", depth), p.Params["Label"], p.Params["Slug"])
		fmt.Printf(" &%p", p.Parent)
		fmt.Println()

		// printMap(p.Params["Label"], p.Params)
		if len(p.Pages) > 0 {
			p.PrintTree(depth + 1)
		}
	}
}
