package mango

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Constants that detected by "gometalinter (goconst)"
const (
	_Yes = "Yes"
	_No  = "No"
	_Md  = ".md"
)

// Application - mango application
type Application struct {
	// sync.RWMutex

	// Absoluted domain path with scheme
	// https://example.loc
	Domain string

	// Absolute path to where binary is
	// config file ".mango" must be there
	binPath string

	// Absolute path to content (folders with .md files)
	ContentPath string

	// Absolute path to web accessible files
	PublicPath string

	// Page tree
	Pages PageList

	// Easy count overall pages and detect duplicates
	// map[Slug]Page
	slugPages PageMap

	// Collectables - pages that are collected
	// Example: "Tag: dog, cat, mouse" --> every tag will point to one *Page
	// Case sensitive
	collections map[string]*Collection

	// Translations
	// If no need to create new .md file but need translate one string
	// translations[lv][Hello] = "Labdien!"
	translations map[string]map[string]string

	// URLTemplates - url templates for pages
	URLTemplates map[string]string

	// Link to server
	Server *Server

	// channel to limit access to App
	chBusy chan bool
}

// NewApplication - create/init new application
// Must be executed only one time
func NewApplication() (*Application, error) {
	app := &Application{}

	// throughput: 1
	// Only one can be manipulating with Application at one moment
	// avoiding concurrency errors
	app.chBusy = make(chan bool, 1)

	// Set defaults
	app.setBinPath()
	app.ContentPath = app.binPath + "/content"
	app.PublicPath = app.binPath + "/public"

	// Default url templates
	// Use {Param} with any Page param
	// TODO: make as string var, not map ?
	app.URLTemplates = map[string]string{
		"Page": "/{Lang}/{Slug}",
		"File": "/{File}",
		// "Collection": "/{collection}/{key}", // /tag/my-tag , /category/Dogs
		// "Group": "/{Lang}/{Slug:[a-z0-9\\-]+}",
	}

	// Configure app by default config file ".mango"
	// Override defaults (as last action)
	app.loadConfig(".mango")

	// Load
	app.LoadContent()

	return app, nil
}

// Detect bin path from where binary executed
func (app *Application) setBinPath() error {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil || strings.Index(path, "/tmp/") >= 0 {
		// if user run "go run *.go" binary will be created in temp folder
		if _path, err := filepath.Abs("."); err == nil {
			path = _path
		}
	}

	app.binPath = path

	return nil
}

// BinPath - get bin path
func (app *Application) BinPath() string {
	return app.binPath
}

// loadConfig using given config filename
// usually ".mango"
// Should not be tested for parallel because used only once in init
func (app *Application) loadConfig(fname string) {
	fpath := app.binPath + "/" + fname
	params := fileToParams(fpath)

	// Overwrite only allowed params

	// Domain which used in sitemap.xml and constructing absolute page url's
	if domain := params["Domain"]; domain != "" {
		domain = strings.TrimPrefix(domain, "http://") // remove default

		if strings.Index(domain, "https://") != 0 {
			domain = "http://" + domain
		}

		// page URL's starts with slash. So skip in domain
		domain = strings.TrimSuffix(domain, "/")
		app.Domain = domain
	}

	// Where content can be found
	if path := params["ContentPath"]; path != "" {
		path, _ = filepath.Abs(path)
		app.ContentPath = filepath.Clean(path)
	}

	// Where web accessible files goes
	if path := params["PublicPath"]; path != "" {
		path, _ = filepath.Abs(path)
		app.PublicPath = filepath.Clean(path)
	}
	os.MkdirAll(app.PublicPath+"/images/", 0755) // where all images from content path moved
	os.MkdirAll(app.PublicPath+"/data/", 0755)   // other file smoved here

	// Template for construction page url's
	if urlTemplate := params["PageURL"]; urlTemplate != "" {
		// Slug must be very specific
		app.URLTemplates["Page"] = urlTemplate
	}
	app.URLTemplates["Page"] = strings.Replace(app.URLTemplates["Page"], "{Slug}", "{Slug:[a-z0-9\\-]+}", -1)

	// Template for construction file url's
	if urlTemplate := params["FileURL"]; urlTemplate != "" {
		app.URLTemplates["File"] = urlTemplate
	}
	app.URLTemplates["File"] = strings.Replace(app.URLTemplates["File"], "{File}", "{File:.+}", -1)

	// Init collections
	// Collections: Tags, Categories, Keywords--> init 3 collection page maps
	app.collections = make(map[string]*Collection, 0) // make anyways
	if params["Collections"] == "" {
		params["Collections"] = "Tags, Categories, Keywords" // default collections
	}
	if ckeys := strings.Split(params["Collections"], ","); len(ckeys) > 0 {
		for _, ckey := range ckeys {
			ckey = strings.TrimSpace(ckey)

			if ckey == "" {
				continue
			}

			// Init
			// Add empty to later know what we are collecting (in app.LoadContent)
			app.collections[ckey] = NewCollection()
		}
	}

}

// LoadContent - Load files to application
// Can be executed more than once times (on .reload file creation)
func (app *Application) LoadContent() {
	app.chBusy <- true // thread-safe

	// Init pageList
	app.slugPages.MakeEmpty()

	// Clear collections
	for ckey := range app.collections {
		app.collections[ckey].MakeEmpty()
	}

	// Page tree
	app.Pages = app.loadPages(app.ContentPath)

	// Post-load operations
	// Edit page after all pages loaded
	app.afterLoadContent()

	// Load translations from every language folder
	app.loadTranslations()

	// Create sitemap.xml under public path
	app.createSitemap()

	<-app.chBusy
}

// Directory to page tree
// TODO: separate goroutine for every language directory listing ?
func (app *Application) loadPages(fpath string) PageList {

	// Collect all pages
	var pages PageList
	if files, dirErr := ioutil.ReadDir(fpath); dirErr == nil {
		for _, f2 := range files {
			if f2.Name()[0] == '.' {
				// Skip config files (e.g. .dir, .defaults..)
				continue
			}

			// Not dir and not .md
			// move to public path
			ext := filepath.Ext(f2.Name())
			ext = strings.ToLower(ext)
			if !f2.IsDir() && ext != _Md {
				images := ".png, .gif, .jpg, .jpeg, .svg," // comma-ended
				mvPath := app.PublicPath + "/images/" + f2.Name()
				if strings.Index(images, ext) == -1 {
					// Images move to /public/data/
					mvPath = app.PublicPath + "/data/" + f2.Name()
				}
				os.Rename(fpath+"/"+f2.Name(), mvPath)
				continue
			}

			// Parse valid page file
			p := app.FileToPage(fpath + "/" + f2.Name())

			if !p.IsYes("IsVisible") {
				// Only visible pages are added
				continue
			}

			// Can't be duplicate slugs
			if p.isDuplicate() {
				p.avoidDuplicate()
			}

			// Add to linear slugPages
			// app.slugPages[p.Params["Slug"]] = p
			app.slugPages.Add(p.Get("Slug"), p)

			// If page is unlisted do not add it to tree
			// (but leave in linear list of pages)
			// That means it can be found by slug,
			// but can't be found among parent childrens
			if !p.IsYes("IsUnlisted") {
				// Add to pageTree
				pages = append(pages, p)
			}

			// After all go deeper.
			// Depth loader must be executed last so top pages are added first
			// Like, "content/en" and "content/lv" are saved and depth pages
			// can reference to them immediately
			//
			// Load sub-pages if it's directory
			if p.IsDir() {
				p.Pages = app.loadPages(p.Get("Path"))
				p.Pages.Sort(p.Get("Sort"))

				// Add parent to received pages
				for _, p2 := range p.Pages {
					// Set Parent page for all sub-pages
					p2.Parent = p
				}
			}

		}
	}

	return pages
}

// Post-load operations
// Edit page after all pages loaded
// For example param: ContentFrom, can be used only after all pages loaded
func (app *Application) afterLoadContent() {

	// Do filter walk but don't collect pages
	app.slugPages.Filter(func(p *Page) bool {

		// What separator to use by appending content
		sepTemplate := []byte("\n{{ Content }}")
		if _sepTemplate := p.Get("ContentTemplate"); _sepTemplate != "" {
			sepTemplate = []byte(_sepTemplate)
		}

		// *** ContentFrom:
		if cfrom := p.Get("ContentFrom"); cfrom != "" {
			if strings.Index(cfrom, ":") > 0 {
				// From collection
				arr := strings.SplitN(cfrom, ":", 2)
				ckey := arr[0]
				citem := arr[1]

				pages := app.CollectionPages(ckey, citem)
				pages.Sort(p.Get("Sort"))

				// Load content from sub-pages
				for _, p3 := range pages {
					content := bytes.Replace(sepTemplate, []byte("{{ Content }}"), p3.Content(), 1)
					p.AppendContent(content)
				}

			} else if page2 := app.Page(cfrom); page2 != nil {
				// From slug

				// Make content
				if page2.IsDir() {

					// Load content from sub-pages
					for _, p3 := range page2.Pages {
						content := bytes.Replace(sepTemplate, []byte("{{ Content }}"), p3.Content(), 1)
						p.AppendContent(content)
					}

				} else {
					// Content from one page
					content := bytes.Replace(sepTemplate, []byte("{{ Content }}"), page2.Content(), 1)
					p.AppendContent(content)
				}

				p.Set("HaveContent", _Yes)
			}

		}

		// *** BreadCrumbs:
		// Add breadcrumb by walking to top by parents
		// For breadcrumbs cant use only filepath because slugs can be different
		p.WalkTop(func(parent *Page) {
			crumbs := parent.Get("Slug") + " / " + p.Get("BreadCrumbs")
			p.Set("BreadCrumbs", crumbs)
		})

		// *** Redirect:
		// 1. Try to get page by redirect slug (if not language root page)
		// 2. add / at the beginning if not absolute url
		if s := p.Get("Redirect"); s != "" {
			url := s
			if p2 := app.Page(s); p2 != nil && p2.IsSet("Level") {
				// Destination page found assign it's URL
				url = p2.Get("URL")

			} else if s[0] != '/' && s[0] != '?' && strings.Index(s, ":") == -1 {
				url = "/" + s
			}
			p.Set("Redirect", url)
		}

		return false
	})
}

// Load translations from every language folder
func (app *Application) loadTranslations() {
	app.translations = make(map[string]map[string]string, 0)

	// Loop only first level (it's language folders)
	for _, p := range app.Pages {
		fpath := p.Get("Path") + "/.translations"
		buf, _ := ioutil.ReadFile(fpath)

		translations := bufToParams(buf, false)
		app.translations[p.Get("Slug")] = translations
	}
}

// Search pages from given top page
func (app *Application) Search(pageSlug, sterm string) PageList {
	page := app.Page(pageSlug) // from where to start search

	// If search is activated when app.LoadContent is in process
	// there could be situation where app.Page(slug) is empty
	if page == nil {
		return nil
	}

	// Success results
	return page.Search(sterm)
}

// CreateSitemap - sitemap.xml
func (app *Application) createSitemap() {
	filepath := app.PublicPath + "/sitemap.xml"
	// tCreated := time.Now()

	contents := `<?xml version="1.0" encoding="UTF-8"?>
<urlset
xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" >
`
	// Language roots
	if strings.Index(app.URLTemplates["Page"], "{Lang") >= 0 {
		for _, p := range app.Pages {
			contents += "<url>\n"
			contents += "\t<loc>" + p.AbsoluteURL() + "</loc>\n"
			contents += "</url>\n"
		}
	}

	// All pageList
	for _, p := range app.slugPages.m {

		if p.IsYes("IsUnlisted") {
			// Unlisted pages must no be in sitemap
			continue
		}

		if p.IsEqual("IsSitemap", "No") {
			// Some pages can be removed only from sitemap
			continue
		}

		if p.IsSet("Redirect") {
			// Redirect pages also not inside
			continue
		}

		contents += "<url>\n"
		contents += "\t<loc>" + p.AbsoluteURL() + "</loc>\n"
		contents += "\t<lastmod>" + p.ModTime().Format(time.RFC3339) + "</lastmod>\n"
		contents += "</url>\n"
	}

	contents += "</urlset>"

	// write whole the body
	ioutil.WriteFile(filepath, []byte(contents), 0644)

}

// IsValidLang - is given language is valid in App scope
func (app *Application) IsValidLang(lang string) bool {
	_, isValid := app.translations[lang]
	return isValid
}

// Print - output app highlights
func (app *Application) Print() {
	log.Println(". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .")
	log.Printf("%20s: %s", "ContentPath", app.ContentPath)
	log.Printf("%20s: %s", "PublicPath", app.PublicPath)
	log.Printf("%20s: %d", "Page count", app.PageCount())
	log.Printf("%20s: %d", "Page (dir) count", len(app.slugPages.Filter(func(p *Page) bool { return p.IsDir() })))
	log.Printf("%20s: %d", "Page (.md) count", len(app.slugPages.Filter(func(p *Page) bool { return !p.IsDir() })))
	log.Printf("%20s: %d", "Collections", len(app.collections))
	log.Println(". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .")
	log.Println()

	// Print every language folder tree
	for _, p := range app.Pages {
		p.PrintTree(0)
	}
	log.Println()

	// Print linear pages by slugs
	app.slugPages.Print()

	// Print linear pages by slugs
	for ckey, c := range app.collections {
		c.Print(ckey)
		log.Println()
	}

}
