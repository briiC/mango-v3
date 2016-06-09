package mango

import (
	"os"
	"strings"
)

// NewPage already linked to app
func (app *Application) NewPage(lang, label string) *Page {
	page := newPage(label)

	// Reload all conent for app if "".reload" file is created in bin path
	reloadFpath := app.BinPath() + "/.reload"
	if _, err := os.Stat(reloadFpath); err == nil {
		// log.Println("[.reload] Reload all pages")
		os.Remove(reloadFpath)
		app.LoadContent() // Reload
		// time.Sleep(time.Second * 1)
		// time.Sleep(time.Millisecond * 1)
	}

	// Set initial language (mandatory)
	page.SetLang(lang)

	// Link page to app
	app.linkPage(page)
	return page
}

// FileToPage for application
func (app *Application) FileToPage(fpath string) *Page {
	page := fileToPage(fpath)
	app.linkPage(page)
	return page
}

// Page - get one page by given slug.
// Slug must be equal and is case-sensitive
func (app *Application) Page(slug string) *Page {
	return app.slugPages.Get(slug)
}

// Assign page to application
// and add some app related params
func (app *Application) linkPage(page *Page) {
	page.App = app

	// Add more params from absolute path
	page.setPathParams()

	// Set Default language on linking page
	// It's validate existing lang
	page.SetLang(page.Get("Lang"))

	// Load page defaults from language root page
	if pDef := app.Page("." + page.Get("Lang") + "-defaults"); pDef != nil {
		page.MergeParams(pDef.params) // fill empty params with defaults
		// After merge check Title
		// Title is not merged of it's special status
		// So we merge it here as exception
		if page.IsEqual("Title", "") {
			page.Set("Title", pDef.Get("Title"))
		}
	}

	// Add "URL" param
	// Only if all other params is set
	if page.ParamsLen() > 0 {
		url := app.URLTemplates["Page"]
		url = page.PopulateParams(url)
		url = "/" + strings.TrimLeft(url, "/") // Fix broken url "//slug/" to "/slug"
		page.Set("URL", url)
	}

	// Add to collections
	for ckey := range app.collections {
		// Is page have such collection key
		if page.IsSet(ckey) {
			// Get this page valuesfrom from c.key
			arr := page.Split(ckey, ",")

			// Every [value: *Page] added to collection by c.key
			for _, itemKey := range arr {
				// Add to app.collections[ckey][itemKey]-> [page1, page2, ..]
				app.collections[ckey].Append(itemKey, page)
			}
		}
	}

	// Reset content because that function holds content normalizer
	// adds correct image and data urls
	page.SetContent(page.Content())

}

// PageCount - total count of pages
func (app *Application) PageCount() int {
	return app.slugPages.Len()
}
