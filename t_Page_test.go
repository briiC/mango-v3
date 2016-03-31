package mango

import "testing"

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
	if !pages[0].IsEqual("Slug", "top-menu") ||
		!pages[1].IsEqual("Slug", "sports") ||
		!pages[2].IsEqual("Slug", "weather") ||
		!pages[3].IsEqual("Slug", "where-is-waldo") {

		pages.Print()
		t.Fatal("Incorrect SearchByParam results")
	}
}
