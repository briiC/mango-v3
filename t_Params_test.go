package mango

import (
	"fmt"
	"testing"
)

func Test_Params(t *testing.T) {
	app, _ := NewApplication()

	// Define filenames and expected result map
	cases := map[string]map[string]string{
		"birthday": nil,
		"golf":     nil,
		"curling":  nil,
		"hello": map[string]string{
			"FileName":    "1_Hello.md",
			"Title":       "Hello",
			"Label":       "Hello",
			"Slug":        "hello",
			"Keywords":    "top, keywords, hello, markdown", // top defaults
			"Author":      "Mango",                          // top defaults
			"Lang":        "en",
			"Ext":         ".md",
			"IsVisible":   "Yes",
			"Level":       "3",
			"HaveContent": "Yes",
			"SortNr":      "1",
			"VisibleFrom": "1420070400000000000", // past
			"VisibleTo":   "{invalid format}",
			"Spaced key":  "",
			"Spaced":      "",
			"BreadCrumbs": "en / en-top-menu / news / ",
		},
		"hello-2": map[string]string{
			"FileName":    "6_Hello.md",
			"Title":       "Hello",
			"Label":       "Hello",
			"IsVisible":   "Yes",
			"Level":       "3",
			"HaveContent": "Yes",
			"SortNr":      "6",
			"VisibleTo":   "4076179200000000000", // late in future
		},
		"cat": map[string]string{
			"Tags":       "animal, pet",
			"Categories": "Housepets",
			"SortNr":     "",
			"Icon":       "cat.png",
			"Keywords":   "top, keywords", // top defaults
			"Author":     "Mango",         // top defaults
		},
		"about-cats": map[string]string{
			"CONTENT": "<h1>Here are info about cats</h1>\n\n<p>Miau!</p>\n",
		},
		"dog": map[string]string{
			"Icon": "animal.png",
		},
		"mango": map[string]string{
			"SortNr":     "1",
			"Tags":       "",
			"Categories": "",
		},
		"sports": map[string]string{
			"IsDir":       "Yes",
			"Level":       "2",
			"Sort":        "Random",
			"Ext":         ".dir",
			"HaveContent": "No",
			"Icon":        "sports.png",
			"Keywords":    "fun, play", // override top defaults
			"Author":      "Sportsman", // top defaults
		},
		"soccer": map[string]string{
			"Icon": "ball.png",
		},
		"fruits": map[string]string{
			"ContentFrom": "fruits",
			"IsDir":       "Yes",
			"HaveContent": "Yes",
		},
	}

	// Run test cases and verify result params
	for slug, cParams := range cases {
		p := app.Page(slug)

		if p != nil && cParams == nil {
			t.Fatal("Page [" + slug + "] must NOT be found")
		} else if p == nil && cParams != nil {
			t.Fatal("Page [" + slug + "] must be found")
		}

		for ckey, cval := range cParams {
			if ckey == "CONTENT" {
				if string(p.Content()) != cval {
					fmt.Printf("\n\n[%s] [%d]\n\n", p.Content(), len(p.Content()))
					t.Fatal(slug, "expected CONTENT: \""+cval+"\"", len(cval))
				}
			} else {
				if !p.IsEqual(ckey, cval) {
					p.Print()
					t.Fatal(slug, "\""+ckey+"\" expected to be \""+cval+"\" -- (Found: ["+p.Get(ckey)+"])")
				}
			}
		}
	}
}
