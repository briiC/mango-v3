package mango

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Setup (before any test has run)

	retCode := m.Run() // run every test

	// Tear down (after all tests done)
	os.RemoveAll("test-files/public/")
	// os.Remove("test-files/public/sitemap.xml")

	os.Exit(retCode)
}

func Test_NewApplication(t *testing.T) {
	// .mango ------------------------------------------------------
	app, err := NewApplication()
	// app.Print()

	if err != nil {
		t.Fatal(app, err)
	}

	// Check default paths
	if app.binPath == "" {
		t.Fatal("binPath is empty")
	}

	// Trim from binPath end to content path end
	if app.Domain != "https://example.loc" {
		t.Fatal("Incorrect app domain", app.Domain)
	}

	// Trim from binPath end to content path end
	if !strings.HasSuffix(app.ContentPath, "/test-files/content") {
		t.Fatal("Default ContentPath must end with /test-files/content. Found: " + app.ContentPath)
	}
	// Trim from binPath end to public path end
	if !strings.HasSuffix(app.PublicPath, "/test-files/public") {
		t.Fatal("Default PublicPath must end with /test-files/public. Found: " + app.PublicPath)
	}

	// .mango2 ------------------------------------------------------
	// Test empty config file
	app2, _ := NewApplication()
	app = app2
	app.loadConfig(".mango2")

	// Trim from binPath end to content path end
	if app.Domain != "http://example.loc" {
		t.Fatal("Incorrect app domain", app.Domain)
	}

	// .mango-empty ------------------------------------------------------
	// Test empty config file
	app3, _ := NewApplication()
	app = app3
	app.loadConfig(".mango-empty")

	// Check default paths
	if app.binPath == "" {
		t.Fatal("binPath is empty")
	}

	// Trim from binPath end to content path end
	if !strings.HasSuffix(app.ContentPath, "/test-files/content") {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}
	// Trim from binPath end to public path end
	if !strings.HasSuffix(app.PublicPath, "/test-files/public") {
		t.Fatal("Default PublicPath must end with /test-files/public")
	}

	// By default must be 3 collections:
	// Tags, Categories, Keywords
	if len(app.collections) != 3 {
		t.Fatal("Must be 3 default collections")
	}

	// Helper function to check valid app language
	if !app.IsValidLang("lv") {
		t.Fatal("App language can be [lv]")
	}
	if app.IsValidLang("xx") {
		t.Fatal("App language can't be [xx]")
	}

}

func Test_AppPageFuncs(t *testing.T) {
	app, _ := NewApplication()
	app.OnReload = func(app *Application) {
		// Set this func, but do nothing
		// Test coverage purpose
	}
	ioutil.WriteFile(app.BinPath()+"/.reload", nil, 0666)

	// Search
	pages := app.Search("en", "oc") // hOCkey, sOCcer, http://remote.lOC/.., localhost..
	if len(pages) != 4 && len(pages) != 3 {
		// NOTE: There can be different count of results depending
		// where test are executed
		// Because of one result tries to get content from "http://localhost"
		// And if it fails that result not counted (one result less)
		pages.Print()
		t.Fatal("Must be found 4 pages")
	}

	// Search - no such slug, no results
	xpages := app.Search("non-existing-slug", "nope")
	if len(xpages) != 0 {
		t.Fatal("Must be found 0 pages")
	}

	// Check reload file
	ioutil.WriteFile(app.binPath+"/.reload", []byte("..."), 0644)

	// new virtual page
	// linked to app, but no parents
	// not listed anywhere
	p := app.NewPage("en", "Virtual reality!")
	p.Set("Custom", "param")
	if !p.IsYes("IsVirtual") ||
		!p.IsEqual("Label", "Virtual reality!") ||
		!p.IsEqual("VirtualSlug", "virtual-reality") ||
		p.App == nil {
		p.Print()
		t.Fatal("Labeled NewPage")
	}

	// Empty label
	p = app.NewPage("en", "")
	p.Set("Custom", "param")
	if !p.IsYes("IsVirtual") ||
		!p.IsEqual("VirtualSlug", "") ||
		!p.IsEqual("Title", "Default title") ||
		p.App == nil {
		p.Print()
		t.Fatal("Empty labeled NewPage")
	}

	// Invalid language
	p = app.NewPage("xx", "")
	p.Set("Custom", "param")
	if !p.IsYes("IsVirtual") ||
		!p.IsEqual("VirtualSlug", "") ||
		!p.IsEqual("Title", "Default title") ||
		!p.IsEqual("Lang", "en") ||
		p.App == nil {
		p.Print()
		t.Fatal("Invalid language NewPage")
	}

}

func Test_AppCollectionFuncs(t *testing.T) {
	app, _ := NewApplication()

	if count := app.CollectionCount(); count != 3 {
		t.Fatal("Collections: incorrect count. Found:", count)
	}

	if count := app.Collection("Tags").Len(); count != 3 {
		t.Fatal("Tags: incorrect count. Found:", count)
	}

	if count := app.Collection("Categories").Len(); count != 4 {
		t.Fatal("Categories: incorrect count. Found:", count)
	}

	if count := app.CollectionPages("Undefined", "nope").Len(); count != 0 {
		t.Fatal("Undefined collection: incorrect count. Found:", count)
	}
}

func Test_FileURLs(t *testing.T) {
	fid := fmt.Sprintf("f-%d", time.Now().UnixNano())
	paths := map[string]string{
		"test-files/content/1_en/" + fid + ".png":  "images",
		"test-files/content/1_en/" + fid + ".JPeg": "images",
		"test-files/content/1_en/" + fid + ".pdf":  "data",
		"test-files/content/1_en/" + fid + ".Pdf":  "data",
	}

	// Crete temp content files
	for fpath := range paths {
		ioutil.WriteFile(fpath, []byte("content"), 0644)
	}

	app, _ := NewApplication()

	// Test location where files must be located after loading app
	for fpath, dirScope := range paths {
		fname := filepath.Base(fpath)
		destPath := app.PublicPath + "/" + dirScope + "/" + fname

		if finfo, _ := os.Stat(destPath); finfo == nil {
			t.Fatalf("[%s] must be moved to [%s]", fpath, destPath)
		}

		os.Remove(fpath)
		os.Remove(destPath)
	}

	// Test urls in content
	page := app.Page("about")
	expected := `<h2>Image urls</h2>

<p><img src="/images/logo.png" alt="img" /><br />
<img src="/images/logo.png" alt="img" /><br />
<img src="/images/logo.png" alt="img" />
<img src="/images/http.png" alt="img" />
<img src="/http.png" alt="img" /><br />
<img src="/lv/logo.png" alt="img" />
<img src="/data/logo.png" alt="img" />
<img src="http://remote.loc/logo.png" alt="img" /><br />
<img src="https://remote.loc/logo.png" alt="img" /><br />
<img src="ftp://remote.loc/logo.png" alt="img" /></p>

<h2>Data urls</h2>

<p><a href="/data/file.pdf">pdf</a><br />
<a href="/data/file.pdf">pdf</a><br />
<a href="/data/file.pdf">pdf</a>
<a href="/data/http.pdf">pdf</a>
<a href="/http.pdf">pdf</a><br />
<a href="/lv/file.pdf">pdf</a>
<a href="/images/file.pdf">pdf</a>
<a href="http://remote.loc/file.pdf">pdf</a><br />
<a href="https://remote.loc/file.pdf">pdf</a><br />
<a href="ftp://remote.loc/file.pdf">pdf</a></p>
`
	if s := tContent(page); string(s) != expected {
		t.Fatalf("Incorrect content [%s]", s)
	}
}
