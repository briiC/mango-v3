package mango

import (
	"strconv"
	"strings"
)

// Page - page with content and params + sub-pages
type Page struct {
	// Linkt to application
	App *Application

	// Params that describe this page
	Params map[string]string

	// Parent page
	Parent *Page

	// Sub-pages for this page
	Pages []*Page
}

// newPage - create/init new page
func newPage(app *Application, fpath string) *Page {

	page := &Page{
		App:    app,
		Params: fileToParams(fpath),
	}

	// some params from path
	page.pathToParams()

	return page
}

// Get some params from path
func (page *Page) pathToParams() {
	if page.App == nil {
		// App must be linked
		return
	}

	// relative path from app.ContentPath
	rpath := strings.TrimPrefix(page.Params["Path"], page.App.ContentPath)

	// Remove filename
	rpath = strings.TrimSuffix(rpath, page.Params["FileName"])

	// split to parts
	rpath = strings.Trim(rpath, "/")
	arr := strings.Split(rpath, "/")

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

	// 1.  Prefix with parent slug
	if page.isDuplicate() && page.Parent != nil {
		page.Params["Slug"] = page.Parent.Params["Slug"] + "-" + page.Params["Slug"]
	}

	// // 2. Prefix with GroupKey
	// if page.isDuplicate() {
	// 	page.Params["Slug"] = page.Params["GroupKey"] + "-" + page.Params["Slug"]
	// }

	// 3. Prefix with Language
	if page.isDuplicate() {
		page.Params["Slug"] = page.Params["Lang"] + "-" + page.Params["Slug"]
	}

	// 4. Loop until unique
	origSlug := page.Params["Slug"]
	for i := 2; page.isDuplicate(); i++ {
		page.Params["Slug"] = origSlug + "-" + strconv.Itoa(i)
	}

}
