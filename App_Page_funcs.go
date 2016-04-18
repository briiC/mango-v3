package mango

import "strings"

// NewPage already linked to app
func (app *Application) NewPage(label string) *Page {
	page := newPage(label)

	// Set default language
	// Check only depth=0 for language
	if len(app.Pages) >= 1 {
		lang := app.Pages[0].Get("Slug")
		page.Set("Lang", lang)
	}

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

	// Load page defaults from language root
	defaults := fileToParams(app.ContentPath + "/" + page.Get("Lang") + "/.defaults")
	page.MergeParams(defaults) // fill empty params with defaults

	// Add "URL" param
	url := app.URLTemplates["Page"]
	url = page.PopulateParams(url)
	url = "/" + strings.TrimLeft(url, "/") // Fix broken url "//slug/" to "/slug"
	page.Set("URL", url)

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
}

// PageCount - total count of pages
func (app *Application) PageCount() int {
	return app.slugPages.Len()
}
