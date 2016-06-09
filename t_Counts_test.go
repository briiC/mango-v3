package mango

import "testing"

func Test_Counts(t *testing.T) {
	app, _ := NewApplication()

	// Pages
	// including ".en-defaults", ".lv-defaults"
	if count := app.PageCount(); count != 41 {
		app.Print()
		t.Fatal("Page count incorrect. Found:", count)
	}

	// Collections
	if count := len(app.collections); count != 3 {
		app.Print()
		t.Fatal("Collection count incorrect. Found:", count)
	}
}

func Test_TagsCounts(t *testing.T) {
	app, _ := NewApplication()

	// Collections: Tags
	if count := app.collections["Tags"].Len(); count != 3 {
		app.Print()
		t.Fatal("[Tags] count incorrect. Found:", count)
	}

	// Collections: Tags: animal
	if count := len(app.collections["Tags"].Get("animal")); count != 5 {
		app.Print()
		t.Fatal("[Tags: animal] count incorrect. Found:", count)
	}

	// Collections: Tags: pet
	if count := len(app.collections["Tags"].Get("pet")); count != 2 {
		app.Print()
		t.Fatal("[Tags: pet] count incorrect. Found:", count)
	}

	// Collections: Tags: nice
	if count := len(app.collections["Tags"].Get("nice")); count != 1 {
		app.Print()
		t.Fatal("[Tags: nice] count incorrect. Found:", count)
	}
}

func Test_CategoriesCounts(t *testing.T) {
	app, _ := NewApplication()

	// Collections: Categories
	if count := app.collections["Categories"].Len(); count != 4 {
		app.Print()
		t.Fatal("[Categories] count incorrect. Found:", count)
	}

}

func Test_KeywordsCounts(t *testing.T) {
	app, _ := NewApplication()

	// Collections: Keywords
	if count := app.collections["Keywords"].Len(); count != 7 {
		app.Print()
		t.Fatal("[Keywords] count incorrect. Found:", count)
	}

}

func Test_SearchCounts(t *testing.T) {
	app, _ := NewApplication()

	// Search
	results := app.Page("en").Search("oC")
	if count := len(results); count != 3 {
		app.Print()
		results.Print()
		t.Fatal("[Search: oC] count incorrect. Found:", count)
	}

	//Search: by param
	results = app.Page("en").SearchByParam("IsDir", _Yes)
	if count := len(results); count != 7 {
		app.Print()
		results.Print()
		t.Fatal("[Search: IsDir] count incorrect. Found:", count)
	}

	// Search: Custom filter function
	results = app.Page("en").Walk(func(p *Page) bool {
		return !p.IsSet("IsDir") || p.IsEqual("IsDir", "No")
	})
	if count := len(results); count != 25 {
		app.Print()
		results.Print()
		t.Fatal("[Search: NOT IsDir] count incorrect. Found:", count)
	}

}
