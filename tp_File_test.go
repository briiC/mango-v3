package mango

import "testing"

func Benchmark_FileToParams_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fileToParams("test-files/content/en/top-menu/1_Simple.md")
		}
	})
}
