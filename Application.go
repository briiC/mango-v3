package mango

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Constants that detected by "gometalinter (goconst)"
const (
	_Yes = "Yes"
	_No  = "No"
	_Md  = ".md"
)

// Application - mango application
type Application struct {
	sync.RWMutex

	// Absolute path to where binary is
	// config file ".mango" must be there
	BinPath string

	// Absolute path to content (folders with .md files)
	ContentPath string

	// Absolute path to web accessable files
	PublicPath string

	// Page tree
	Pages []*Page

	// Easy count overall pages and detect duplicates
	// map[Slug]Page
	slugPages PageMap

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
	app.ContentPath = app.BinPath + "/content"
	app.PublicPath = app.BinPath + "/public"

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
		if path, err = filepath.Abs("."); err != nil {
			// Can't be tested because this will be only on "go run *.go"
			return err
		}
	}

	app.BinPath = path

	return nil
}

// loadConfig using given config filename
// usually ".mango"
// TODO: concurrent
func (app *Application) loadConfig(fname string) {
	fpath := app.BinPath + "/" + fname
	params := fileToParams(fpath)
	// if len(params) == 0 {
	// 	// log.Printf("Error: Empty OR not exists [%s] config file", fname)
	// }

	// Overwrite only allowed params

	if path := params["ContentPath"]; path != "" {
		path, _ = filepath.Abs(path)
		app.Lock()
		app.ContentPath = path
		app.Unlock()
	}

	if path := params["PublicPath"]; path != "" {
		path, _ = filepath.Abs(path)
		app.Lock()
		app.PublicPath = path
		app.Unlock()
	}
}

// LoadContent - Load files to application
// Can be executed more than once times (on .reload file creation)
func (app *Application) LoadContent() {
	app.chBusy <- true // thread-safe

	// Init pageList
	app.slugPages.MakeEmpty()

	// Page tree
	app.Pages = app.loadPages(app.ContentPath)

	<-app.chBusy
}

// Directory to page tree
// TODO: separate goroutine for every language directory listing ?
func (app *Application) loadPages(fpath string) PageList {

	// Collect all pages
	var pages PageList
	if files, dirErr := ioutil.ReadDir(fpath); dirErr == nil {
		for _, f2 := range files {
			if f2.Name()[:1] == "." {
				// Skip config files (e.g. .dir, .defaults..)
				continue
			}

			p := app.NewPage(fpath + "/" + f2.Name())

			if p.Params["IsVisible"] != _Yes {
				// Only visible pages are added
				continue
			}

			// Can't be duplicate slugs
			if p.isDuplicate() {
				p.avoidDuplicate()
			}

			// Load sub-pages if it's directory
			if p.Params["IsDir"] == _Yes {
				p.Pages = app.loadPages(p.Params["Path"])

				// Sort by default
				if len(p.Pages) >= 2 {
					switch p.Params["Sort"] {
					case "Reverse":
						sort.Sort(sort.Reverse(p.Pages))
					case "Random":
						p.Pages.Randomize()
					default:
						sort.Sort(p.Pages)
					}
				}

				// Add parent to received pages
				for _, p2 := range p.Pages {
					// Set Parent page for all sub-pages
					p2.Parent = p
				}
			}

			// Add to linear slugPages
			// app.slugPages[p.Params["Slug"]] = p
			app.AddPage(p)

			// Add to pageTree
			pages = append(pages, p)

		}
	}

	return pages
}

// NewPage for application
func (app *Application) NewPage(fpath string) *Page {
	return newPage(app, fpath)
}

// Page - get one page by given slug.
// Slug must be equal and is case-sensitive
func (app *Application) Page(slug string) *Page {
	// if len(app.chBusy) > 0 {
	// 	return nil
	// }
	return app.slugPages.Get(slug)
}

// AddPage - adds page to pageList
// Returns (bool) if slug changed to avoid duplicates
func (app *Application) AddPage(page *Page) {

	// Do not overwrite existing slug
	if page.isDuplicate() {
		page.avoidDuplicate()
	}

	app.slugPages.Set(page.Get("Slug"), page)
}

// RemovePage - remove page from app and returns (bool) if succeeded
func (app *Application) RemovePage(slug string) {
	// app.chBusy <- true

	// remove from tree
	page := app.Page(slug)
	if page != nil && page.Parent != nil {

		// This page is under parent pages
		pages := page.Parent.Pages
		for i, p := range pages {
			// Found this page index in slice
			if p.IsEqual("Slug", slug) {

				// remove from slice
				pages = append(pages[:i], pages[i+1:]...)
			}
		}

	}

	// Remove from slugPages
	app.slugPages.Remove(slug)

	// <-app.chBusy
}
