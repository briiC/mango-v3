package mango

import (
	"bytes"
	"fmt"
	"testing"
)

// Parsing datetimes
func Test_PageContent(t *testing.T) {
	app, _ := NewApplication()

	cases := map[string][]byte{
		"lava": []byte("This is very deep file"),
		"golf": []byte("# Golf"),
		"cold": []byte("Winter is coming.."),
		"one-more": []byte("# Header line\n" +
			"\n" +
			"- Some **markdown** syntax.\n" +
			"- And some <b>HTML</b> synta too."),
		// "copy-cat": []byte("# Waldo here"),
	}

	// loop cases
	for slug, expected := range cases {
		content := app.Page(slug).Content
		if !bytes.Equal(content, expected) {
			fmt.Printf("\n\n::: FOUND: %s\n\n", content)
			fmt.Printf("::: EXPECTED: %s\n\n", expected)
			t.Fatal("Invalid content in [", app.Page(slug).Params["Path"], "]")
		}
	}

}

// Parsing datetimes
func Test_PageFuncs(t *testing.T) {
	app, _ := NewApplication()
	page := app.Page("golf")

	// Set/Get param
	page.Set("Label", "Golfing") //was Golf
	if page.Get("Label") != "Golfing" {
		t.Fatal("Label must be changed")
	}

	// Param helpers
	if !page.IsEqual("Label", "Golfing") {
		t.Fatal("ERROR: IsEqual")
	}
	if !page.IsSet("Label") {
		t.Fatal("ERROR: IsSet")
	}
	if !page.IsYes("HaveContent") {
		t.Fatal("ERROR: IsYes")
	}
	if page.IsNo("HaveContent") {
		t.Fatal("ERROR: IsNo")
	}
	if page.IsDir() {
		t.Fatal("ERROR: IsDir: Must not be directory")
	}

	// Search functions
	// Sarch only on "en" language scope
	// Order of search must not change because using slice not map
	pages := app.Page("en").Search("w")
	if !pages[0].IsEqual("Slug", "weather") ||
		!pages[1].IsEqual("Slug", "where-is-waldo") ||
		!pages[2].IsEqual("Slug", "waldo") {

		pages.Print()
		t.Fatal("Incorrect Search results")
	}

	// Search functions
	// Sarch only on "en" language scope
	// Order of search must not change because using slice not map
	pages = app.Page("en").SearchByParam("IsDir", "Yes")
	if !pages[0].IsEqual("Slug", "en-top-menu") ||
		!pages[1].IsEqual("Slug", "sports") ||
		!pages[2].IsEqual("Slug", "weather") ||
		!pages[3].IsEqual("Slug", "where-is-waldo") {

		pages.Print()
		t.Fatal("Incorrect SearchByParam results")
	}

	// Custom filter function
	// Sarch only on "en" language scope
	pages = app.Page("en").Walk(func(p *Page) bool {
		return !p.IsSet("IsDir") || p.IsEqual("IsDir", "No")
	})
	if !pages[0].IsEqual("Slug", "copy-cat") ||
		!pages[1].IsEqual("Slug", "simple-slug-oh") ||
		!pages[2].IsEqual("Slug", "one-more") ||
		!pages[3].IsEqual("Slug", "last-in-line") {

		pages.Print()
		t.Fatal("Incorrect Custom Walk results")
	}

	// Trye to Change Slug with setter
	// Must NOT be changed
	page = app.Page("golf")
	page.Set("Slug", "golfing")

	page = app.Page("golf") // Still can found
	if page == nil {
		t.Fatal("Slug must not be changed")
	}

	page = app.Page("golfing") // Must NOT be found
	if page != nil {
		t.Fatal("New slug must not be found")
	}

}
