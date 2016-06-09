package mango

import "testing"

// Concurrency testing with all client operations
func Benchmark_App_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.LoadContent()

			app.NewPage("en", "Page.md")
			app.FileToPage("test-files/content/1_en/top-menu/1_Simple.md")

			// Page
			app.Page("hello")

			// Search
			app.Search("en", "oc") // hOCkey, sOCcer

			// PageMap
			p := &Page{}
			p.Set("Slug", "slug-x") // slug must be set for  slugPages
			app.slugPages.Add("slug-x", p)
			app.PageCount()
			app.slugPages.Remove("slug-x")

			// Params
			p.Set("Custom", "hi")
			p.Get("Custom")
			p.RemoveParam("Custom")
			p.SetValue("KeyInt", 102)
			p.SetValue("KeyBool", true)
			p.Params()

			// Collection
			app.CollectionCount()
			app.Collection("Tags").Append("tag-x", &Page{})
			app.Collection("Tags").Get("tag-x")
			app.Collection("Tags").Len() //count of items insife Tags
			app.Collection("Tags").Remove("tag-x")

			app.Collection("Categories").Append("cat-x", &Page{})
			app.CollectionPages("Categories", "cat-x")
			app.Collection("Categories").Len() //count of items insife Tags
			app.Collection("Categories").Remove("cat-x")

			// Reload page
			if p := app.Page("monkey"); p != nil {
				p.ReloadContent()
			}

		}
	})
}
