package mango

import (
	"fmt"
	"html/template"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

var (
	// FuncMap - To use in html/template FuncMap
	defaultFuncMap = template.FuncMap{
		"T":         T,
		"Get":       tGet,
		"Content":   tContent,
		"Page":      tPage,
		"HTML":      tHTML,
		"MdToHTML":  tMdToHTML,
		"Slice":     tSlice,
		"SliceFrom": tSliceFrom,
		// "PageURL":        PageURL,
		// "FileURL":        FileURL,
		// "GetParams":      GetParams,
		"ToTags":      tParseToTags,
		"CurrentYear": tCurrentYear,
		"DateFormat":  tDateFormat,
		"FileURL":     tFileURL,
		"Print":       tPrint,
		"Loop":        tLoop,
		"Split":       tSplitToSlice,
	}
)

// T - Translate string to given language
func T(page *Page, s string) string {
	if page.App == nil {
		return s // cant translate w/o app
	}

	if translated, ok := page.App.translations[page.Get("Lang")][s]; ok {
		return translated
	}

	return s
}

// Get param from Page or params map
func tGet(page interface{}, key string) string {

	switch page.(type) {
	case *Page:
		return page.(*Page).Get(key)
	case map[string]string:
		return page.(map[string]string)[key]
	}

	return ""
}

// Get param from Page or params map
func tContent(page *Page) template.HTML {
	return template.HTML(page.Content())
}

// Get Page by given slug
// Give Application context
func tPage(page *Page, slug string) *Page {
	p := page.App.Page(page.Get("Lang") + "-" + slug)
	if p == nil {
		p = page.App.Page(slug)
	}
	return p
}

// Convert given params to HTML
func tHTML(args ...interface{}) template.HTML {
	s := fmt.Sprintf("%s", args...)
	return template.HTML(s)
}

// Markdown To HTML - Parse string (markdown) to html
// used in tmeplates
func tMdToHTML(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

// Slice - used in template to Slice
// from - to
// use it in template as:
//		Slice $Arr 0 4
func tSlice(pages PageList, from, to int) PageList {
	if to >= len(pages) {
		// slice as many can
		to = len(pages)
	}

	if from < 0 || to <= 0 {
		return nil
	}

	return pages[from:to]
}

// SliceFrom - used in template to Slice
// from - (to end)
func tSliceFrom(pages PageList, from int) PageList {
	if from >= len(pages) || from < 0 {
		return nil
	}
	return pages[from:]
}

// ParseToTags - parse to tags for html output
// @line can have multiple tag definitions, but separated with comma ","
func tParseToTags(page *Page, codeLang string, params ...string) template.HTML {
	codeLang = strings.ToLower(codeLang)
	rawLine := strings.Trim(strings.Join(params, ","), " ,\n\t")

	htmlStr := ""

	if rawLine == "" {
		// Try to use page param if no val given
		switch codeLang {
		case "javascript":
			rawLine = page.Get("JavaScript")
		case "css":
			rawLine = page.Get("CSS")
		case "breadcrumbs":
			rawLine = page.Get("BreadCrumbs")
		}
	}

	// To avoid splitting in wrong comma first replace all "\,"
	// Later we replace it back to comma
	rawLine = strings.Replace(rawLine, "\\,", "%%COMMA%%", -1)

	// Split lines to arr
	lines := strings.Split(rawLine, ",")

	for _, line := range lines {
		if line == "" {
			continue
		}

		line = strings.Trim(line, " ")
		isURL := (line[:1] == "/" || line[:4] == "http")

		switch codeLang {
		case "javascript":
			if isURL {
				// src
				htmlStr += "<script type=\"text/javascript\" src=\"" + line + "\"></script>\n"
			} else {
				// plain code
				htmlStr += "<script type=\"text/javascript\">" + line + "</script>\n"
			}
		case "css":
			if isURL {
				// href
				htmlStr += "<link rel=\"stylesheet\" href=\"" + line + "\" type=\"text/css\" />\n"
			} else {
				// plain code
				htmlStr += "<style type=\"text/css\">" + line + "</style>\n"
			}

		case "breadcrumbs":
			line = strings.Trim(line, "/ \t\n")
			slugs := strings.Split(line, "/")
			var arr []string
			for _, slug := range slugs {
				slug = strings.TrimSpace(slug)
				if p := page.App.Page(slug); p != nil {
					if lvl, _ := strconv.Atoi(p.Get("Level")); lvl >= 2 {
						// Skip root levels
						// If need root levels add it in html by yourself
						arr = append(arr, "<a href=\""+p.Get("URL")+"\" title=\""+p.Get("Title")+"\">"+p.Get("Title")+"</a>")
					}
				}
			}
			htmlStr += strings.Join(arr, " / ")
		}
	}

	// Replace native commas (not separators) to real commas
	htmlStr = strings.Replace(htmlStr, "%%COMMA%%", ",", -1)

	return template.HTML(htmlStr)
}

func tCurrentYear() int {
	year := time.Now().Year()
	return year
}

func tFileURL(page *Page, parts ...string) string {
	// get prefix
	arr := strings.SplitN(page.App.URLTemplates["File"], "{File", 2)
	prefix := arr[0]

	// construct based on url file template
	url := prefix + strings.Join(parts, "/")
	return path.Clean(url)
}

func tDateFormat(layout, s string) string {
	if t, err := toTime(s); err == nil {
		return t.Format(layout)
	}

	return s // return as given
}

func tPrint(p interface{}) string {
	switch p.(type) {
	case *Page:
		p.(*Page).Print()
	}
	return "..."
}

// return slice filled with items of given number
func tLoop(s interface{}) []int {

	n := 1 //default

	// Convert any type to int
	if i, err := strconv.Atoi(fmt.Sprintf("%v", s)); err == nil {
		n = i
	}

	// Not returning only empty slice
	// but fill it with actual numbers
	// so in template we get correct number on:
	// {{ range $i := Loop 10 }} {{ $i }} {{ end }}
	// starting with 1 not 0
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1 // slice[0] = 1
	}

	return arr
}

// Split strings by separator
func tSplitToSlice(s, sep string) []string {
	return strings.Split(s, sep)
}
