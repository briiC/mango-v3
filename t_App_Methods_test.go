package mango

import (
	"fmt"
	"testing"
)

// Parsing datetimes
func Test_LoadConfig(t *testing.T) {
	app, _ := NewApplication()

	// App runs .mango by default
	// this is just comfort code as example
	app.loadConfig(".mango")
	if app.ContentPath[len(app.binPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}

	// override with new params
	app.loadConfig("test-files/.mango-empty")
	if app.ContentPath[len(app.binPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}

}

func Test_LoadPages(t *testing.T) {
	app, _ := NewApplication() //auto-load

	// Count
	if app.slugPages.Len() != 27 {
		// pages := app.loadPages(app.ContentPath)
		for _, p := range app.Pages {
			p.PrintTree(0)
		}
		fmt.Println()

		t.Fatalf("Must be exact number of pages. Found %d", app.slugPages.Len())
	}

	// hidden folder must not exists
	if app.Page("hidden") != nil {
		t.Fatalf("Hidden (dot prefixed) folder must not be visible")
	}

	// Check if duplicates are correct
	if app.Page("one-more") == nil ||
		app.Page("one-more-2") == nil ||
		app.Page("one-more-3") == nil ||
		app.Page("one-more-4") == nil {
		t.Fatalf("All wannabe-duplicates must exist with modified slug")
	}

	// Very deep file correct
	// checking foldr
	if app.Page("lava") == nil ||
		app.Page("lava").Params["Level"] != "6" ||
		app.Page("lava").Params["Lang"] != "lv" ||
		app.Page("lava").Params["GroupKey"] != "top-menu" {
		printMap("Lava", app.Page("lava").Params)
		t.Fatal("Lava page OR Lava params not correct")
	}

	// Test DEFAULT order for TopMenu pages
	tmPages := app.Page("en-top-menu").Pages
	if tmPages[0].Params["Slug"] != "simple-slug-oh" ||
		tmPages[1].Params["Slug"] != "one-more" ||
		tmPages[2].Params["Slug"] != "last-in-line" {

		for i, p := range tmPages {
			fmt.Println("\t\t", i, p.Params["SortNr"], p.Params["Slug"])
		}
		fmt.Println()

		t.Fatal("Order of TopMenu pages not correct")
	}

	// Test REVERSE order for Sports pages
	spPages := app.Page("sports").Pages
	if spPages[0].Params["Slug"] != "hockey" ||
		spPages[1].Params["Slug"] != "golf" ||
		spPages[2].Params["Slug"] != "baseball" {

		for i, p := range spPages {
			fmt.Println("\t\t", i, p.Params["SortNr"], p.Params["Slug"])
		}
		fmt.Println()

		t.Fatal("Order of Sports pages not correct")
	}

	// Test RANDOM order
	// pseudo check. If SortNr are set it could be random
	// because by default SortNr are not set there
	wPages := app.Page("where-is-waldo").Pages
	if wPages[0].Params["SortNr"] == "" ||
		wPages[0].Params["SortNr"] == "0" ||
		len(wPages) != 4 {

		for i, p := range wPages {
			fmt.Println("\t\t", i, p.Params["SortNr"], p.Params["Slug"])
		}
		fmt.Println()

		t.Fatal("Order of WALDO must be random")
	}

}

func Test_AppNewPage(t *testing.T) {
	app, _ := NewApplication()

	// Empty label
	p := app.NewPage("")
	// printMap("", p.Params)
	if p.Params["IsVirtual"] != "Yes" ||
		p.Params["VirtualSlug"] != "" ||
		p.App == nil {
		t.Fatal("Empty labeled NewPage")
	}

	// Labeled
	p = app.NewPage("Hello page!")
	// printMap("Hello page!", p.Params)
	if p.Params["IsVirtual"] != "Yes" ||
		p.Params["Label"] != "Hello page!" ||
		p.Params["VirtualSlug"] != "hello-page" ||
		p.App == nil {
		t.Fatal("Labeled NewPage")
	}
}

func Test_AppFileToPage(t *testing.T) {
	app, _ := NewApplication()

	// Empty label
	p := app.FileToPage("")
	if len(p.Params) != 0 ||
		p.App == nil {
		printMap("", p.Params)
		t.Fatal("Empty labeled FileToPage")
	}

	// Labeled
	p = app.FileToPage("Hello page!")
	if len(p.Params) != 0 ||
		p.App == nil {
		printMap("Hello page!", p.Params)
		t.Fatal("Labeled FileToPage")
	}

	// Existing
	p = app.FileToPage("test-files/content/en/top-menu/1_Simple.md")
	if len(p.Params) == 0 ||
		p.Params["Slug"] != "simple-slug-oh" ||
		p.App == nil {
		printMap("Existing", p.Params)
		t.Fatal("Existing path for FileToPage")
	}
}
