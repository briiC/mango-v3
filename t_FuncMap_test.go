package mango

import (
	"testing"
	"time"
)

// TestNew - create and check application
func Test_FuncMap(t *testing.T) {
	app, _ := NewApplication()
	page := app.Page("cat")

	// tT
	if s := tT(page, "Hello"); s != "Labdien" {
		t.Fatalf("Incorrect translation [%s]", s)
	}
	if s := tT(page, "Undefined"); s != "Undefined" {
		t.Fatalf("Incorrect translation [%s]", s)
	}
	p := newPage("Not linked")
	if s := tT(p, "Hello"); s != "Hello" {
		t.Fatalf("Incorrect translation [%s]", s)
	}

	// tGet
	if s := tGet(page, "Label"); s != "Cat" {
		t.Fatalf("Incorrect param [%s]", s)
	}
	params := map[string]string{
		"Label": "From map",
	}
	if s := tGet(params, "Label"); s != "From map" {
		t.Fatalf("Incorrect param [%s]", s)
	}
	if s := tGet(nil, "Label"); s != "" {
		t.Fatalf("Incorrect param [%s]", s)
	}

	// tContent
	if s := tContent(page); s != "<p>Miau!</p>\n" {
		t.Fatalf("Incorrect content [%s]", s)
	}

	// tPage
	if p := tPage(page, "cat"); p == nil {
		t.Fatal("Page must be found")
	}
	if p := tPage(page, "undefined"); p != nil {
		t.Fatal("Page must NOT be found")
	}

	if s := tHTML("<b>html</b>"); s != "<b>html</b>" {
		t.Fatalf("Incorrect html [%s]", s)
	}

	if s := tMdToHTML("**markdown**"); s != "<p><strong>markdown</strong></p>\n" {
		t.Fatalf("Incorrect markdown [%s]", s)
	}

	pages := page.Parent.Pages // cat -> animals

	// tSlice
	if _pages := tSlice(pages, 0, 1); len(_pages) != 1 {
		t.Fatalf("Incorrect slice [%d]", len(_pages))
	}
	if _pages := tSlice(page.Pages, 0, 1); len(_pages) != 0 {
		t.Fatalf("Incorrect slice [%d]", len(_pages))
	}

	if _pages := tSliceFrom(pages, 1); len(_pages) != 4 {
		t.Fatalf("Incorrect slice [%d]", len(_pages))
	}
	if _pages := tSliceFrom(page.Pages, 1); len(_pages) != 0 {
		t.Fatalf("Incorrect slice [%d]", len(_pages))
	}

	// tParseToTags
	if s := tParseToTags(page, "javascript", "a.js,, /js/b.js"); s != "<script type=\"text/javascript\">a.js</script>\n<script type=\"text/javascript\" src=\"/js/b.js\"></script>\n" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}
	if s := tParseToTags(page, "css", "a.css,, /css/b.css"); s != "<style type=\"text/css\">a.css</style>\n<link rel=\"stylesheet\" href=\"/css/b.css\" type=\"text/css\" />\n" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}
	if s := tParseToTags(page, "css", ""); s != "" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}
	if s := tParseToTags(page, "javascript", ""); s != "" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}
	if s := tParseToTags(page, "breadcrumbs", ""); s != "<a href=\"/en/animals.html\" title=\"Animals\">Animals</a>" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}

	// tLoop
	if arr := tLoop(3); len(arr) != 3 || arr[0] != 1 {
		t.Fatalf("Incorrect slice for loop [%v]", arr)
	}
	if arr := tLoop("2"); len(arr) != 2 || arr[0] != 1 {
		t.Fatalf("Incorrect slice for loop [%v]", arr)
	}

	//
	if d := tCurrentYear(); d != time.Now().Year() {
		t.Fatalf("Incorrect year [%d]", d)
	}

	if s := tFileURL(page, "logo.png"); s != "/logo.png" {
		t.Fatalf("Incorrect file url [%s]", s)
	}

	// Datetimes
	if s := tDateFormat("02.01.2006", "1984-07-02"); s != "02.07.1984" {
		t.Fatalf("Incorrect datetime parse [%s]", s)
	}
	if s := tDateFormat("02.01.2006", "1461584775277491501"); s != "25.04.2016" {
		t.Fatalf("Incorrect datetime parse [%s]", s)
	}
	if s := tDateFormat("02.01.2006", "xxx"); s != "xxx" {
		t.Fatalf("Incorrect datetime parse [%s]", s)
	}

}
