package mango

import "testing"

func Benchmark_PageParamsRead_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			slug := app.Page("golf").Get("Slug")
			_ = slug[1:] // only to avoid warning of unused variable
		}
	})
}

func Benchmark_PageParamsGetSet_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Get some page
			page := app.Page("en").Search("w")[0]

			// unsafe
			// app.pageList["golf"].Params["Label"] = "write"

			//safe
			// get/set
			page.Set("Label", app.Page("cold").Get("Label"))
			page.Set("Slug", "new-golf")
		}
	})
}
