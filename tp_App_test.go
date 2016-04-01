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

func Benchmark_AppLoadPage_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.FileToPage("Page.md")
		}
	})
}

func Benchmark_AppMixed_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.LoadContent()
			app.Page("golf")
			app.FileToPage("Page.md")
		}
	})
}
