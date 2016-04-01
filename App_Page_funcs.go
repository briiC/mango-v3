package mango

// NewPage already linked to app
func (app *Application) NewPage(label string) *Page {
	page := newPage(label)
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
	page.setPathParams()

	// Add to collections
	for ckey := range app.collections {
		// Is page have such collection key
		if page.IsSet(ckey) {
			// Get this page valuesfrom from c.key
			arr := page.Split(ckey, ",")

			// Every [value: *Page] added to collection by c.key
			for _, val := range arr {
				// fmt.Println(page.Get("Slug"), "\t", val)
				// Add to app.collections[ckey][val] []*Page
				// 	      app.collections[Tag][my-tag] []*Page
				app.collections[ckey][val] = append(app.collections[ckey][val], page)
			}
		}
	}
}
