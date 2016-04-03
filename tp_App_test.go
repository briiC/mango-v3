package mango

import "testing"

// Concurrency testing with all client operations
func Benchmark_App_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.LoadContent()

			app.NewPage("Page.md")
			app.FileToPage("test-files/content/en/top-menu/1_Simple.md")

			// Page
			app.Page("hello")

			// Search
			app.Search("en", "oc") // hOCkey, sOCcer

			// PageMap
			p := &Page{}
			p.Set("Slug", "slug-x") // slug must be set for  slugPages
			app.slugPages.Add("slug-x", p)
			app.slugPages.Len()
			app.slugPages.Remove("slug-x")

			// Collection
			app.collections["Tags"].Append("tag-x", &Page{})
			app.collections["Tags"].Len()
			app.collections["Tags"].Remove("tag-x")

		}
	})
}
