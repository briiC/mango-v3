package mango

import (
	"fmt"
	"strconv"
	"strings"
)

// Page - page with content and params + sub-pages
type Page struct {
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

// Check if page is duplicate slug
func (page *Page) isDuplicate() bool {
	_, isDuplicate := page.App.pageList[page.Params["Slug"]]
	return isDuplicate
}

// Generate unique slug based on old one
func (page *Page) avoidDuplicate() {
	// Suffix loop by count until unique
	origSlug := page.Params["Slug"]
	for i := 2; page.isDuplicate(); i++ {
		page.Params["Slug"] = origSlug + "-" + strconv.Itoa(i)
	}

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
