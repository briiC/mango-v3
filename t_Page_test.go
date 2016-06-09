package mango

import (
	"fmt"
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
	if !page.IsNegation("XXX") {
		t.Fatal("ERROR: IsNegation")
	}
	if page.IsNegation("HaveContent") {
		t.Fatal("ERROR: IsNegation")
	}
	if page.IsDir() {
		t.Fatal("ERROR: IsDir: Must not be directory")
	}

	page.Set("My", "Custom")
	page.RemoveParam("My")
	if page.IsSet("My") {
		t.Fatal("ERROR: RemoveParam")
	}

	if paramCount := len(page.Params()); paramCount != 21 {
		t.Fatal("ERROR: Params(): Found:", paramCount)
	}

	page.SetValue("IntVal", 102)
	if !page.IsEqual("IntVal", "102") {
		t.Fatal("int value not correct")
	}

	page.SetValue("BoolVal", true)
	if !page.IsEqual("BoolVal", "Yes") {
		t.Fatal("bool value not correct")
	}

	page.SetValue("BoolVal", false)
	if !page.IsEqual("BoolVal", "No") {
		t.Fatal("bool value not correct")
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
	p := app.NewPage("en", "Virtual")
	p.Set("Path", "../no-such-file")
	if p.ReloadContent() != false {
		fmt.Println(p.Get("Path"))
		t.Fatalf("File doesn't exists to reload")
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

	// Set page lang
	if page.SetLang("lv"); page.Get("Lang") != "lv" {
		t.Fatal("Page lang must be [lv]. Found: " + page.Get("Lang"))
	}
	if page.SetLang("en"); page.Get("Lang") != "en" {
		t.Fatal("Page lang must be [en]")
	}
	if page.SetLang("xx"); page.Get("Lang") == "xx" {
		t.Fatal("Page lang can't be [xx]")
	}

}
