package mango

import "testing"

// Parsing datetimes
func Test_PageFuncs(t *testing.T) {
	app, _ := NewApplication()
	page := app.Page("hello")
	// page.Print()

	// Set/Get param
	page.Set("Label", "Hello again!") //was Golf
	if page.Get("Label") != "Hello again!" {
		t.Fatal("Label must be changed")
	}

	// Param helpers
	if !page.IsEqual("Label", "Hello again!") {
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

	// Trye to Change Slug with setter
	// Must NOT be changed
	page = app.Page("hello")
	page.Set("Slug", "goodbye")

	page = app.Page("hello") // Still can found
	if page == nil {
		t.Fatal("Slug must not be changed")
	}

	page = app.Page("goodbye") // Must NOT be found
	if page != nil {
		t.Fatal("New slug must not be found")
	}
}
