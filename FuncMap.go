package mango

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

var (
	// FuncMap - To use in html/template FuncMap
	defaultFuncMap = template.FuncMap{
		"T":         tT,
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
	}
)

// Translate string to given language
func tT(page *Page, s string) string {
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
	p := page.App.Page(slug)
	if p == nil {
		p = page.App.Page(page.Get("Lang") + "-" + slug)
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
	if to >= len(pages) || from < 0 || to <= 0 {
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
func tParseToTags(codeLang string, params ...string) template.HTML {
	codeLang = strings.ToLower(codeLang)
	rawLine := params[0]
	htmlStr := ""

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
			break

		case "css":
			if isURL {
				// href
				htmlStr += "<link rel=\"stylesheet\" href=\"" + line + "\" type=\"text/css\" />\n"
			} else {
				// plain code
				htmlStr += "<style type=\"text/css\">" + line + "</style>\n"
			}
			break

			// case "breadcrumbs":
			// 	if len(params) >= 2 {
			// 		separator := App.BreadcrumbSeparator
			// 		if separator == "" {
			// 			separator = " / " //default separator
			// 		}
			// 		crumbs := strings.Split(line, separator)
			// 		urls := strings.Split(params[1], " ")
			// 		for index, crumbFilename := range crumbs {
			// 			crumbFilename = strings.Trim(crumbFilename, " ")
			// 			if crumbFilename != "" {
			// 				label, _, _, _, _, _, _ := FilenameToParams(crumbFilename)
			// 				if label == "_" || label == "-" || label[:1] == "." {
			// 					continue
			// 				}
			// 				htmlStr += "<a href=\"" + urls[index] + "\">" + label + "</a> "
			// 				htmlStr += separator
			// 			}
			// 		}
			// 		htmlStr = strings.TrimSuffix(htmlStr, separator)
			// 	}
			// 	break
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
