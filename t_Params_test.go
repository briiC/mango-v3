package mango

import (
	"fmt"
	"testing"
)

// TestNew - create and check application
func Test_MergeParams(t *testing.T) {
	params := map[string]string{
		"A": "aaa",
		"B": "bbb",
		"C": "ccc",
	}
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

	params = mergeParams(params, params2, params3)
	if len(params) != 5 {
		fmt.Printf("Merged: %+v", params)
		t.Fatal(params, "Expected too see 5 params. Found only ", len(params))
	}
}
