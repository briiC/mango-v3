package mango

import "testing"

func Benchmark_MergeParams(b *testing.B) {
	params := map[string]string{} //global

	params2 := map[string]string{
		"B": "bbb",
		"C": "ccc",
		"D": "ddd",
	}
	params3 := map[string]string{
		"C": "ccc",
		"D": "ddd",
		"E": "eee",
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			params = map[string]string{
				"A": "aaa",
				"B": "bbb",
				"C": "ccc",
			}
			params = mergeParams(params, params2, params3)
		}
	})
}
