package mango

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Constants that detected by "gometalinter (goconst)"
const (
	_Yes = "Yes"
	_No  = "No"
	_Md  = ".md"
)

// Application - mango application
type Application struct {
	isBusy bool

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
	// Page slice (not tree)
	// map[Slug]Page
	pageList map[string]*Page
}

// NewApplication - create/init new application
func NewApplication() (*Application, error) {
	app := &Application{}

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
func (app *Application) loadConfig(fname string) {
	// Busy while loading config
	app.isBusy = true
	defer func() { app.isBusy = false }()

	fpath := app.BinPath + "/" + fname
	params := fileToParams(fpath)
	// if len(params) == 0 {
	// 	// log.Printf("Error: Empty OR not exists [%s] config file", fname)
	// }

	// Overwrite only allowed params

	if path := params["ContentPath"]; path != "" {
		path, _ = filepath.Abs(path)
		app.ContentPath = path
	}

	if path := params["PublicPath"]; path != "" {
		path, _ = filepath.Abs(path)
		app.PublicPath = path
	}
}

// LoadContent - Load files to application
func (app *Application) LoadContent() {
	// Busy while loading files
	app.isBusy = true
	defer func() { app.isBusy = false }()

	// Init pageList
	app.pageList = make(map[string]*Page, 0)

	// Page tree
	app.Pages = app.loadPages(app.ContentPath)

}

// NewPage for application
func (app *Application) NewPage(fpath string) *Page {
	return newPage(app, fpath)
}

// Directory to page tree
func (app *Application) loadPages(fpath string) []*Page {

	// Get info about fpath
	// Only dir can be used for loading pages
	f, fErr := os.Stat(fpath)
	if fErr != nil || !f.IsDir() {
		return nil
	}

	// Collect all pages
	var pages []*Page
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

			if p.Params["IsDir"] == _Yes {
				// Load sub-pages if it's directory
				p.Pages = app.loadPages(p.Params["Path"])
				for _, p2 := range p.Pages {
					// Set Parent page for all sub-pages
					p2.Parent = p
				}
			}

			// Add to linear pageList
			app.pageList[p.Params["Slug"]] = p

			// Add to pageTree
			pages = append(pages, p)

		}
	}

	return pages
}
