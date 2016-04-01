package mango

import (
	"io/ioutil"
	"log"
	"testing"
)

// Do test coverage for output functions
// but do not output anywhere. Silence.
func Test_CoverOutputFuncs(t *testing.T) {
	app, _ := NewApplication()
	pages := app.Page("en").Search("w")

	// *** Disable log outputs
	log.SetOutput(ioutil.Discard)

	// App output
	app.Print()

	// PageMap output
	app.slugPages.Print()

	// PageList output
	pages.Print()

	// *Page
	app.Page("en").PrintTree(0)

	// map[string]string
	printMap(app.Page("en").Get("Slug"), app.Page("en").Params)

}
