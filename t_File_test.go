package mango

import "testing"

// TestNew - create and check application
func Test_FilenameToParams(t *testing.T) {
	// Define filenames and expected result map
	cases := map[string]map[string]string{
		".mango": {
			"Ext":       ".mango",
			"Label":     "",
			"Slug":      "",
			"IsVisible": "No",
		},
		"Simple.md": {
			"Ext":       ".md",
			"Label":     "Simple",
			"Title":     "Simple",
			"Slug":      "simple",
			"IsVisible": "Yes",
		},
		"path.to/some/file/Simple.md": {
			"Ext":       ".md",
			"Label":     "Simple",
			"Title":     "Simple",
			"Slug":      "simple",
			"IsVisible": "Yes",
		},
		"path.to/some/file/ŪTF 8.md/": {
			"Ext":       ".md",
			"Label":     "ŪTF 8",
			"Title":     "ŪTF 8",
			"Slug":      "utf-8",
			"IsVisible": "Yes",
		},
		"65_With sort number.md": {
			"Ext":       ".md",
			"Label":     "With sort number",
			"Title":     "With sort number",
			"Slug":      "with-sort-number",
			"IsVisible": "Yes",
			"SortNr":    "65",
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

// // TestNew - create and check application
// func Test_FileToParams(t *testing.T) {
// 	// Define filenames and expected result map
// 	cases := map[string]map[string]string{
// 		"test-files/.params": map[string]string{
// 			"Ext":           ".params",
// 			"Label":         "",
// 			"Slug":          "",
// 			"B":             "bbb",
// 			"Valid":         "counts only first \":\"",
// 			"StillValid":    "ignoring prespace",
// 			"A":             "rewriting key A",
// 			"NotMultiLine":  "This not part of \"MultiLine\" param",
// 			"NotMultiLine2": "have \\ but not at the end",
// 			"ModTime":       "BEGIN: 14", // 1458839061048723130
// 		},
// 		"not-exists": map[string]string{
// 			"Ext":       "",
// 			"Label":     "",
// 			"Slug":      "",
// 			"IsVisible": "",
// 		},
// 		"test-files/content/en/top-menu/.defaults": map[string]string{
// 			"Icon":      "default.ico",
// 			"IsVisible": "No",
// 		},
// 		"test-files/content/en/top-menu/1_Simple.md": map[string]string{
// 			"Icon":     "default.ico", // this param comes from .defaults
// 			"SubIcon":  "subicon.ico", // this param comes from .subdefaults ^1
// 			"DeepIcon": "deep.ico",    // this param comes from .subdefaults ^2
// 			"Over":     "by in-file",  //
// 			"Sub":      "first",       //
//
// 			"IsVisible": "Yes",
// 			"Slug":      "simple-slug-oh",
// 			"Ext":       ".md",
// 			"SortNr":    "1",
// 			"Label":     "Simple changed",
// 			"Title":     "Simple changed",
//
// 			"Keywords":    "A, B, C, D",
// 			"Path":        "test-files/content/en/top-menu/1_Simple.md", //relative?
// 			"VisibleFrom": "1426896000000000000",
// 			"VisibleTo":   "4077877497000000000",
// 		},
// 		"test-files/content/en/top-menu/2_One more.md": map[string]string{
// 			"Label":       "One more",
// 			"Slug":        "one-more",
// 			"Keywords":    "A, B, C",
// 			"Title":       "Title is changed",
// 			"VisibleFrom": "invalid:date:time",
// 		},
// 		"test-files/content/en/top-menu/Weather/Cold.md": map[string]string{
// 			"Label": "Cold",
// 			"Slug":  "cold",
// 			"Title": "Cold",
// 		},
// 		"test-files/content/en/top-menu/Sports": map[string]string{
// 			"Label":       "Sports",
// 			"Slug":        "sports",
// 			"Title":       "Sports",
// 			"IsDir":       "Yes",
// 			"IsVisible":   "Yes",
// 			"HaveContent": "No",
// 			"Ext":         ".dir",
// 			"Path":        "ENDS: /Sports",
//
// 			"Icon":     "default.ico", // this param comes from .defaults
// 			"SubIcon":  "subicon.ico", // this param comes from .subdefaults ^1
// 			"DeepIcon": "deep.ico",    // this param comes from .subdefaults ^2
// 		},
// 		"test-files/content/en/top-menu/Weather": map[string]string{
// 			"Label":       "Weather",
// 			"Slug":        "weather",
// 			"Title":       "Weather",
// 			"IsDir":       "Yes",
// 			"IsVisible":   "Yes",
// 			"HaveContent": "No",
// 			"Ext":         ".dir",
// 			"Path":        "ENDS: /Weather",
//
// 			"Icon":     "snow.ico",    // this param comes from .defaults
// 			"SubIcon":  "subicon.ico", // this param comes from .subdefaults ^1
// 			"DeepIcon": "deep.ico",    // this param comes from .subdefaults ^2
// 		},
// 		"test-files/content/en/top-menu/Sports/2_Baseball.md": map[string]string{
// 			"IsVisible":   "Yes",
// 			"VisibleFrom": "1426896000000000000",
// 			"VisibleTo":   "", // or not set
// 		},
// 		"test-files/content/en/top-menu/Sports/22_Hockey.md": map[string]string{
// 			"IsVisible":   "Yes",
// 			"VisibleFrom": "", // or not set
// 			"VisibleTo":   "4077734400000000000",
// 		},
// 		"test-files/content/en/top-menu/Sports/021_Soccer.md": map[string]string{
// 			"IsVisible":   "No",
// 			"VisibleFrom": "4077734400000000000",
// 			"VisibleTo":   "",   // or not set
// 			"SortNr":      "21", // not 021
// 		},
// 	}
//
// 	// Run test cases and verify result params
// 	for filename, cParams := range cases {
// 		fParams := fileToParams(filename)
// 		for ckey, cval := range cParams {
// 			notValid := false
//
// 			// Test if value STARTS with correct
// 			if strings.Index(cval, "BEGIN:") == 0 {
// 				beginVal := strings.TrimSpace(cval[6:]) // len(BEGIN:) == 6
// 				if strings.Index(fParams[ckey], beginVal) != 0 {
// 					notValid = true
// 				} else {
// 					continue
// 				}
// 			}
//
// 			// Test if value ENDS with correct
// 			if strings.Index(cval, "ENDS:") == 0 {
// 				endVal := strings.TrimSpace(cval[5:]) // len(BEGIN:) == 6
// 				_cval := fParams[ckey][len(fParams[ckey])-len(endVal):]
// 				if _cval != endVal {
// 					notValid = true
// 				} else {
// 					continue
// 				}
// 			}
//
// 			// Test exact match
// 			if fParams[ckey] != cval {
// 				notValid = true
// 			}
//
// 			// If not valid print map and show error
// 			if notValid {
// 				printMap(filename, fParams)
// 				t.Fatal(filename, "\""+ckey+"\" expected to be \""+cval+"\" -- (Found: ["+fParams[ckey]+"])")
// 			}
// 		}
// 	}
//
// }
