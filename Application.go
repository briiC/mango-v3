package mango

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	// sync.RWMutex

	// Absolute path to where binary is
	// config file ".mango" must be there
	binPath string

	// Absolute path to content (folders with .md files)
	ContentPath string

	// Absolute path to web accessable files
	PublicPath string

	// Page tree
	Pages []*Page

	// Easy count overall pages and detect duplicates
	// map[Slug]Page
	slugPages PageMap

	// Collectables - pages that are collected
	// Example: "Tag: dog, cat, mouse" --> every tag will point to one *Page
	// Case sensitive
	collections map[string]*Collection

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

	app.binPath = path

	return nil
}

// loadConfig using given config filename
// usually ".mango"
// Should not be tested for parallel because used only once in init
func (app *Application) loadConfig(fname string) {
	fpath := app.binPath + "/" + fname
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

			p := app.FileToPage(fpath + "/" + f2.Name())

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
			app.slugPages.Add(p.Get("Slug"), p)

			// If page is unlisted do not add it to tree
			// (but leave in linear list of pages)
			// That means it can be found by slug,
			// but can't be found among parent childrens
			if !p.IsYes("IsUnlisted") {
				// Add to pageTree
				pages = append(pages, p)
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

		// Content from slug
		if slug := p.Get("ContentFrom"); slug != "" {
			if page2 := app.Page(slug); page2 != nil {

				// What seperator to use by appending content
				sepTemplate := []byte("\n{{ Content }}")
				if _sepTemplate := p.Get("ContentTemplate"); _sepTemplate != "" {
					sepTemplate = []byte(_sepTemplate)
				}

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

				p.Set("HaveContent", "Yes")

			}
		}

		return false
	})
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

// Print - output app highlights
func (app *Application) Print() {
	log.Println(". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .")
	log.Printf("%20s: %s", "ContentPath", app.ContentPath)
	log.Printf("%20s: %s", "PublicPath", app.PublicPath)
	log.Printf("%20s: %d", "Page count", app.slugPages.Len())
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
	log.Println()

	// Print linear pages by slugs
	for ckey, c := range app.collections {
		c.Print(ckey)
		log.Println()
	}

}
