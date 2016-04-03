package mango

import "testing"

// Parsing datetimes
func Test_PageMap(t *testing.T) {
	pm := NewPageMap()

	// Add +3
	pm.Add("slug-x", &Page{})
	pm.Add("slug-y", &Page{})
	pm.Add("slug-z", &Page{})
	if pm.Len() != 3 {
		t.Fatal("Pages must be added")
	}

	// Get
	if page := pm.Get("slug-x"); page == nil {
		t.Fatal("Page must be found")
	}
	if page := pm.Get("slug-y"); page == nil {
		t.Fatal("Page must be found")
	}
	if page := pm.Get("slug-z"); page == nil {
		t.Fatal("Page must be found")
	}

	//TODO: pm.Filter

	// Remove -2
	pm.Remove("slug-y")
	pm.Remove("slug-z")
	if pm.Len() != 1 {
		t.Fatal("Pages must be removed")
	}

	// Clear
	pm.MakeEmpty()
	if pm.Len() != 0 {
		t.Fatal("All pages must be cleared")
	}
}
