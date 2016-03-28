package mango

import (
	"os"
	"path/filepath"
	"strings"
)

// Application - mango application
type Application struct {

	// Absolute path to where binary is
	// config file ".mango" must be there
	BinPath string

	// Absolute path to content (folders with .md files)
	ContentPath string

	// Absolute path to web accessable files
	PublicPath string
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
