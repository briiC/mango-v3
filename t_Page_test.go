package mango

import (
	"strings"
	"testing"
	"time"
)

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

	// Reload content (page)
	page = app.Page("my-secret-post")
	page.Set("ModTime", "0") // simulate file content changed
	content := strings.TrimSpace(string(page.Content()))
	if content != "<p>Auto reloaded.</p>" {
		t.Fatalf("Content not correct. Must be reloaded. Found: [%s]", content)
	}

	// Reload content (dir)
	page = app.Page("fruits")
	page.Set("ModTime", "0") // simulate file content changed
	if page.ReloadContent() != false {
		t.Fatalf("Directories can not be reloaded")
	}

	// Urls
	if url := page.AbsoluteURL(); url != "https://example.loc/en/fruits.html" {
		t.Fatalf("Incorrect absolute url: %s", url)
	}

	// ModTime
	page.Set("ModTime", "xxx") // simulate file content changed
	if dt := page.ModTime(); dt.Format("01.02.2006") != time.Now().Format("01.02.2006") {
		t.Fatalf("Incorrect modtime: %v", dt)
	}
	page.Set("ModTime", "xxxxxxxxxxxxxx") // simulate file content changed
	if dt := page.ModTime(); dt.Format("01.02.2006") != time.Now().Format("01.02.2006") {
		t.Fatalf("Incorrect modtime: %v", dt)
	}

	// Paging
	page.Paging(2, 1, 2)
	isValid := page.IsEqual("PPrev", "1") &&
		page.IsEqual("PNum", "2") &&
		page.IsEqual("PNext", "0") &&
		page.IsEqual("PSize", "1") &&
		page.IsEqual("PFrom", "1") &&
		page.IsEqual("PTo", "2") &&
		page.IsEqual("PTotalPages", "2") &&
		page.IsEqual("PTotalItems", "2")
	if !isValid {
		page.Print()
		t.Fatalf("Incorrect paging 1")
	}

	page = app.Page("cat")
	page.Paging(0, 0, 0)
	isValid = page.Get("PPrev") == "0" &&
		page.IsEqual("PNum", "1") &&
		page.IsEqual("PNext", "0") &&
		page.IsEqual("PSize", "1") &&
		page.IsEqual("PFrom", "0") &&
		page.IsEqual("PTo", "0") &&
		page.IsEqual("PTotalPages", "0") &&
		page.IsEqual("PTotalItems", "0")
	if !isValid {
		page.Print()
		t.Fatalf("Incorrect paging 2")
	}

	page = app.Page("cat")
	page.Paging(99, 99, 99)
	isValid = "0" == page.Get("PPrev") &&
		page.IsEqual("PNum", "1") &&
		page.IsEqual("PNext", "0") &&
		page.IsEqual("PSize", "99") &&
		page.IsEqual("PFrom", "0") &&
		page.IsEqual("PTo", "0") &&
		page.IsEqual("PTotalPages", "0") &&
		page.IsEqual("PTotalItems", "0")
	if !isValid {
		page.Print()
		t.Fatalf("Incorrect paging 3")
	}

}
