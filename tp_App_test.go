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

func Benchmark_AppLoadConfig_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.loadConfig(".mango")
		}
	})
}

func Benchmark_AppAddPage_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.AddPage(app.Page("golf")) // ugly slugs created
		}
	})
}

func Benchmark_AppRemovePage_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			page := app.Page("golf")
			app.AddPage(page) // ugly slugs created
			app.RemovePage(page.Get("Slug"))

		}
	})
}

func Benchmark_AppNewPage_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.NewPage("Page.md")
		}
	})
}

func Benchmark_AppMixed_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.LoadContent()
			app.Page("golf")
			app.NewPage("Page.md")
		}
	})
}
