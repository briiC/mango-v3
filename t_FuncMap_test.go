package mango

import (
	"testing"
	"time"
)

// TestNew - create and check application
func Test_FuncMap(t *testing.T) {
	app, _ := NewApplication()
	page := app.Page("cat")

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

	if s := tContent(page); s != "<p>Miau!</p>\n" {
		t.Fatalf("Incorrect content [%s]", s)
	}

	if p := tPage(page, "cat"); p == nil {
		t.Fatal("Page must be found")
	}

	if s := tHTML("<b>html</b>"); s != "<b>html</b>" {
		t.Fatalf("Incorrect html [%s]", s)
	}

	if s := tMdToHTML("**markdown**"); s != "<p><strong>markdown</strong></p>\n" {
		t.Fatalf("Incorrect markdown [%s]", s)
	}

	pages := page.Parent.Pages // cat -> animals

	if _pages := tSlice(pages, 0, 1); len(_pages) != 1 {
		t.Fatalf("Incorrect slice [%d]", len(_pages))
	}

	if _pages := tSliceFrom(pages, 1); len(_pages) != 4 {
		t.Fatalf("Incorrect slice [%d]", len(_pages))
	}

	if s := tParseToTags("javascript", "a.js, /js/b.js"); s != "<script type=\"text/javascript\">a.js</script>\n<script type=\"text/javascript\" src=\"/js/b.js\"></script>\n" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}

	if s := tParseToTags("css", "a.css, /css/b.css"); s != "<style type=\"text/css\">a.css</style>\n<link rel=\"stylesheet\" href=\"/css/b.css\" type=\"text/css\" />\n" {
		t.Fatalf("Incorrect parse to tags [%s]", s)
	}

	if d := tCurrentYear(); d != time.Now().Year() {
		t.Fatalf("Incorrect year [%d]", d)
	}

	if s := tFileURL(page, "logo.png"); s != "/logo.png" {
		t.Fatalf("Incorrect file url [%s]", s)
	}
}
