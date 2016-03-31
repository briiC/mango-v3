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

func Benchmark_AppLoadConfig_Parallel(b *testing.B) {
	app, _ := NewApplication()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			app.loadConfig(".mango")
		}
	})
}
