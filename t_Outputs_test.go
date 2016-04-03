package mango

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// Do test coverage for output functions
// but do not output anywhere. Silence.
func Test_CoverOutputFuncs(t *testing.T) {
	app, _ := NewApplication()
	pages := app.Page("en").Search("w")

	// *** Disable log outputs
	log.SetOutput(ioutil.Discard)

	// enable it after tests done, to not discard logs in other tests
	defer log.SetOutput(os.Stdout)

	// App output
	app.Print()

	// PageMap output
	app.slugPages.Print()

	// PageList output
	pages.Print()

	// *Page
	app.Page("en").Print()
	app.Page("en").PrintTree(0)

	// map[string]string
	printMap(app.Page("en").Get("Slug"), app.Page("en").Params)

}
