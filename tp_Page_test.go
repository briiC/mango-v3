package mango

import "testing"

func Benchmark_PageParamsRead_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			slug := app.pageList["golf"].Params["Slug"]
			_ = slug[1:] // only to avoid warning of unused variable
		}
	})
}

func Benchmark_PageParamsGetSet_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// unsafe
			// app.pageList["golf"].Params["Slug"] = "write"

			//safe
			// get/set
			app.pageList["golf"].Set("Slug", app.pageList["cold"].Get("Slug"))
		}
	})
}
