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
	enPage := app.pageList["en"]
	pages := app.pageList["en"].Search("w")

	// Disable log outputs
	log.SetOutput(ioutil.Discard)

	// PageList output
	pages.Print()

	// *Page
	enPage.PrintTree(0)
}
