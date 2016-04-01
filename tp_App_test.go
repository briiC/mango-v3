package mango

import "testing"

func Benchmark_AppLoadContent_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.LoadContent()
		}
	})
}

func Benchmark_AppPage_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.Page("golf")
		}
	})
}

func Benchmark_AppVirtualPage_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.NewPage("Page.md")
			app.FileToPage("test-files/content/en/top-menu/1_Simple.md")
		}
	})
}

func Benchmark_AppMixed_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.LoadContent()

			app.NewPage("Page.md")
			app.FileToPage("test-files/content/en/top-menu/1_Simple.md")

			app.Page("golf")

			// Add
			app.slugPages.Add("slug-x", &Page{})

			// Remove
			app.slugPages.Remove("slug-x")

		}
	})
}
