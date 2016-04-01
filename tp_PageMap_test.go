package mango

import "testing"

func Benchmark_PageMapOperations_Parallel(b *testing.B) {
	app, _ := NewApplication()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			// Add
			app.slugPages.Add(&Page{
				Params: map[string]string{"Slug": "slug-x"},
			})

			// Get
			app.slugPages.Get("slug-x")

			// Remove
			app.slugPages.Remove("slug-x")

			// Filter
			app.slugPages.Filter(func(p *Page) bool { return p.IsDir() })
		}
	})
}
