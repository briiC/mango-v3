package mango

import (
	"fmt"
	"testing"
)

func printMap(fname string, m map[string]string) {
	fmt.Println("@ ---", fname, "---")
	for key, val := range m {
		fmt.Printf("%20s: %s \n", key, val)
	}
	fmt.Println()
}

// TestNew - create and check application
func Test_FilenameToParams(t *testing.T) {
	// Define filenames and expected result map
	cases := map[string]map[string]string{
		".mango": map[string]string{
			"Ext":       ".mango",
			"Label":     "",
			"Slug":      "",
			"IsVisible": "No",
		},
		"Simple.md": map[string]string{
			"Ext":       ".md",
			"Label":     "Simple",
			"Slug":      "simple",
			"IsVisible": "Yes",
		},
		"path.to/some/file/Simple.md": map[string]string{
			"Ext":       ".md",
			"Label":     "Simple",
			"Slug":      "simple",
			"IsVisible": "Yes",
		},
		"path.to/some/file/ŪTF 8.md/": map[string]string{
			"Ext":       ".md",
			"Label":     "ŪTF 8",
			"Slug":      "utf-8",
			"IsVisible": "Yes",
		},
		"65_With sort number.md": map[string]string{
			"Ext":       ".md",
			"Label":     "With sort number",
			"Slug":      "with-sort-number",
			"IsVisible": "Yes",
			"SortNr":    "65",
		},
		"1_01.01.2000-26.05.2001_Date range.md": map[string]string{
			"Ext":       ".md",
			"Label":     "Date range",
			"Slug":      "date-range",
			"IsVisible": "No",
			"SortNr":    "1",
			"DateFrom":  "2000-01-01 00:00:00 +0000 UTC",
			"DateTo":    "2001-05-26 23:59:00 +0000 UTC",
		},
		"1_01.01.2000-26.05.2099_Date range active.md": map[string]string{
			"Ext":       ".md",
			"Label":     "Date range active",
			"Slug":      "date-range-active",
			"IsVisible": "Yes",
			"SortNr":    "1",
			"DateFrom":  "2000-01-01 00:00:00 +0000 UTC",
			"DateTo":    "2099-05-26 23:59:00 +0000 UTC",
		},
	}

	// Run test cases and verify result params
	for filename, cParams := range cases {
		fParams := filenameToParams(filename)
		for ckey, cval := range cParams {
			if fParams[ckey] != cval {
				printMap(filename, fParams)
				t.Fatal(filename, "\""+ckey+"\" expected to be \""+cval+"\" -- (Found: ["+fParams[ckey]+"])")
			}
		}
	}
}

// TestNew - create and check application
func Test_FileToParams(t *testing.T) {
	// Define filenames and expected result map
	cases := map[string]map[string]string{
		".mango": map[string]string{
			"Ext":         ".mango",
			"Label":       "",
			"Slug":        "",
			"A":           "aaa",
			"B":           "bbb",
			"~StillValid": "mandatory only \":\"",
		},
		"not-exists": map[string]string{
			"Ext":       "",
			"Label":     "",
			"Slug":      "",
			"IsVisible": "",
		},
		"test-files/content/en/top-menu/.defaults": map[string]string{
			"Icon":      "default.ico",
			"IsVisible": "No",
		},
		"test-files/content/en/top-menu/1_Simple.md": map[string]string{
			"Icon":      "default.ico", // this param comes from .defaults
			"SubIcon":   "subicon.ico", // this param comes from .subdefaults ^1
			"DeepIcon":  "deep.ico",    // this param comes from .subdefaults ^2
			"IsVisible": "Yes",
			"Slug":      "simple-slug-oh",
			"Ext":       ".md",
			"SortNr":    "1",
			"Label":     "Simple changed",
		},
	}

	// Run test cases and verify result params
	for filename, cParams := range cases {
		fParams := fileToParams(filename)
		for ckey, cval := range cParams {
			if fParams[ckey] != cval {
				printMap(filename, fParams)
				t.Fatal(filename, "\""+ckey+"\" expected to be \""+cval+"\" -- (Found: ["+fParams[ckey]+"])")
			}
		}
	}

}
