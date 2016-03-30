package mango

import (
	"bytes"
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
	if len(app.pageList) != 26 {
		// pages := app.loadPages(app.ContentPath)
		for _, p := range app.Pages {
			p.PrintTree(0)
		}
		fmt.Println()

		t.Fatalf("Must be exact number of pages. Found %d", len(app.pageList))
	}

	// hidden folder must not exists
	if app.pageList["hidden"] != nil {
		t.Fatalf("Hidden (dot prefixed) folder must not be visible")
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

	// Test DEFAULT order for TopMenu pages
	tmPages := app.pageList["top-menu"].Pages
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
	spPages := app.pageList["sports"].Pages
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
	wPages := app.pageList["where-is-waldo"].Pages
	if wPages[0].Params["SortNr"] == "" ||
		wPages[0].Params["SortNr"] == "0" ||
		len(wPages) != 4 {

		for i, p := range wPages {
			fmt.Println("\t\t", i, p.Params["SortNr"], p.Params["Slug"])
		}
		fmt.Println()

		t.Fatal("Order of WALDO must be random")
	}

	// pages := app.loadPages(app.ContentPath)
	for _, p := range app.Pages {
		p.PrintTree(0)
	}
	fmt.Println()

	//
	for slug, p := range app.pageList {
		fmt.Printf("%20s\t", slug)
		// fmt.Printf("P %-16p ", p)
		// fmt.Printf("Pr %-16p ", p.Parent)
		fmt.Printf("%-5s ", p.Params["IsDir"])
		fmt.Printf("L %-3s ", p.Params["Level"])
		fmt.Printf("%5d bytes ", len(p.Content))
		fmt.Println()
	}
	fmt.Println()

}

func Test_Content(t *testing.T) {
	app, _ := NewApplication() //auto-load

	cases := map[string][]byte{
		"lava": []byte("This is very deep file"),
		"golf": []byte("# Golf"),
		"cold": []byte("Winter is coming.."),
		"one-more": []byte("# Header line\n" +
			"\n" +
			"- Some **markdown** syntax.\n" +
			"- And some <b>HTML</b> synta too."),
	}

	// loop cases
	for slug, expected := range cases {
		content := app.pageList[slug].Content
		if !bytes.Equal(content, expected) {
			fmt.Printf("\n\n::: FOUND: %s\n\n", content)
			fmt.Printf("::: EXPECTED: %s\n\n", expected)
			t.Fatal("Invalid content in [", app.pageList[slug].Params["Path"], "]")
		}
	}
}
