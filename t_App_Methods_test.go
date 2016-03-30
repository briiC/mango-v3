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
	if app.ContentPath[len(app.BinPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}

	// override with new params
	app.loadConfig("test-files/.mango-empty")
	if app.ContentPath[len(app.BinPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}

}

func Test_LoadPages(t *testing.T) {
	app, _ := NewApplication() //auto-load

	// Count
	if len(app.pageList) != 19 {
		t.Fatalf("Must be exact number of pages. Found %d", len(app.pageList))
	}

	// Check if duplicates are correct
	if app.pageList["one-more"] == nil ||
		app.pageList["one-more-2"] == nil ||
		app.pageList["one-more-3"] == nil ||
		app.pageList["one-more-4"] == nil {
		t.Fatalf("All wannabe-duplicates must exist with modified slug")
	}

	// Very deep file correct
	// checking foldr
	if app.pageList["lava"] == nil ||
		app.pageList["lava"].Params["Level"] != "6" ||
		app.pageList["lava"].Params["Lang"] != "lv" ||
		app.pageList["lava"].Params["GroupKey"] != "left-menu" {
		printMap("Lava", app.pageList["lava"].Params)
		t.Fatal("Lava page OR Lava params not correct")
	}

	// pages := app.loadPages(app.ContentPath)
	for _, p := range app.Pages {
		p.PrintTree(0)
	}
	fmt.Println()

	//
	for slug, p := range app.pageList {
		fmt.Printf("%20s &%p &%p %s\n", slug, p, p.Parent, p.Params["Level"])
		// printMap("xx", p.Params)
	}

}
