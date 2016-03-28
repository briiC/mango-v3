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
	app, _ := NewApplication()

	// pages := app.loadPages(app.ContentPath)
	for _, p := range app.Pages {
		p.Print(0)
	}
	fmt.Println()

	//
	for slug, p := range app.pageList {
		fmt.Printf("%20s &%p &%p\n", slug, p, p.Parent)
		// printMap("xx", p.Params)
	}

}
