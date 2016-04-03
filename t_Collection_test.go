package mango

import "testing"

// Parsing datetimes
func Test_Collection(t *testing.T) {
	col := NewCollection()

	// Append +3
	col.Append("tag-x", &Page{})
	col.Append("tag-x", &Page{})
	col.Append("tag-x", &Page{})

	// Append +1
	col.Append("tag-y", &Page{})

	// Append +2
	col.Append("tag-z", &Page{})
	col.Append("tag-y", &Page{})

	if col.Len() != 3 {
		t.Fatal("Pages must be added")
	}

	// Get
	if pages := col.Get("tag-x"); len(pages) != 3 {
		t.Fatal("Incorrect count:", len(pages))
	}
	if pages := col.Get("tag-y"); len(pages) != 2 {
		t.Fatal("Incorrect count:", len(pages))
	}
	if pages := col.Get("tag-z"); len(pages) != 1 {
		t.Fatal("Incorrect count:", len(pages))
	}

	//TODO: col.Filter

	// Remove -2
	col.Remove("tag-y")
	col.Remove("tag-z")
	if col.Len() != 1 {
		t.Fatal("Pages must be removed")
	}

	// Clear
	col.MakeEmpty()
	if col.Len() != 0 {
		col.Print("collection test")
		t.Fatal("All pages must be cleared")
	}
}
