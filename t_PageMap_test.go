package mango

import "testing"

// Parsing datetimes
func Test_PageMap(t *testing.T) {
	app, _ := NewApplication()

	count := app.slugPages.Len()

	// Add
	app.slugPages.Add(&Page{
		Params: map[string]string{"Slug": "slug-x"},
	})
	if count+1 != app.slugPages.Len() {
		t.Fatal("Page must be added")
	}

	// Get
	if page := app.slugPages.Get("slug-x"); !page.IsEqual("Slug", "slug-x") {
		t.Fatal("Page must be found")
	}

	// Remove
	app.slugPages.Remove("slug-x")
	if count != app.slugPages.Len() {
		t.Fatal("Page must be removed")
	}

	// Clear
	app.slugPages.MakeEmpty()
	if app.slugPages.Len() != 0 {
		t.Fatal("All pages must be cleared")
	}
}
