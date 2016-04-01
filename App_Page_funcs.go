package mango

// NewPage already linked to app
func (app *Application) NewPage(label string) *Page {
	page := newPage(label)
	page.linkToApp(app)
	return page
}

// FileToPage for application
func (app *Application) FileToPage(fpath string) *Page {
	page := fileToPage(fpath)
	page.linkToApp(app)
	return page
}

// Page - get one page by given slug.
// Slug must be equal and is case-sensitive
func (app *Application) Page(slug string) *Page {
	return app.slugPages.Get(slug)
}
