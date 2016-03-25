package mango

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Read file contents and convert to map of params
// We do not check for error because we only want to get/or not file
// Importance order:
// 1. Content params
// 2. Filename params
// 3. (same depth) .defaults
// 4 -n. (up n depth) .subdefaults
// 5. .mango config file params
func fileToParams(fpath string) map[string]string {
	// get given filepath directory path
	pwd, _ := filepath.Abs(filepath.Dir(fpath))

	// Check real file
	finfo, fErr := os.Stat(fpath)
	if os.IsNotExist(fErr) {
		// Not exists
		// fileToParams expect to file be created
		// otherwise params empty
		return map[string]string{}
	}

	// Is it directory
	if fErr == nil && finfo.IsDir() {
		// Special file to find
		fpath += "/.dir"
	}

	// Get params from filename
	// Parsing first, but assigning afer content params
	params2 := filenameToParams(fpath)

	// Set existing file system file params
	params2["Path"] = fpath
	if finfo.IsDir() {
		// Get original Path
		params2["Path"] = fpath[:len(fpath)-5] // trim .dir
	}
	params2["ModTime"] = fmt.Sprint(finfo.ModTime().UnixNano()) //TODO: or strconv faster?

	// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {

	// Raw file contents
	// Not checking for read error. If can't read empty content
	buf, _ := ioutil.ReadFile(fpath)
	var bufHeader, bufContent []byte

	// Split raw buf to variables
	sep := []byte("\n+++") // it's no problem to leave \n if front of content
	if params2["Ext"] != ".md" {
		// Not .md file, so use as params file
		bufHeader = buf

	} else if bytes.Index(buf, sep) >= 0 {
		// have header and content separated by content separator
		arr := bytes.SplitN(buf, sep, 2)
		bufHeader = arr[0]
		bufContent = arr[1]
	} else {
		// If not splited use raw buf as content
		bufContent = buf
	}

	// ** Params
	params := bufToParams(bufHeader) // first assign what we can from buf

	// ** Content
	bufContent = bytes.TrimSpace(bufContent)
	if bufContent == nil {
		params["HaveContent"] = "No"
	} else {
		params["HaveContent"] = "Yes"
	}

	// ** Load extra file params
	params3 := make(map[string]string)

	// only for content files or directories
	if params2["Ext"] == ".md" || finfo.IsDir() {

		// Same depth .defaults
		params3 = fileToParams(pwd + "/.defaults")

		// Up level .subdefaults (until can't found)
		subfilepath := pwd
		// can't merge with params3 already, because param order will break
		subparams := make(map[string]string, 0)
	SUB:
		subfilepath, _ = filepath.Abs(subfilepath + "/../") //one up
		_subparams := fileToParams(subfilepath + "/.subdefaults")
		if len(_subparams) > 0 {
			subparams = mergeParams(subparams, _subparams)
			goto SUB
		}
		params3 = mergeParams(params3, subparams)

	}

	// Title is tricky. We need special treatment.
	// Title is based on Label if empty (always set)
	if params["Title"] == "" && params["Label"] != "" {
		params["Title"] = params["Label"]
	}

	// ** Merge params correctly. First param map is more important
	// file <---- filename <---- defaults <- subdefaults
	params = mergeParams(params, params2, params3)

	return params
}

// Parse given bytes to map of params
func bufToParams(buf []byte) map[string]string {
	params := make(map[string]string, 0)

	nl := []byte("\n") // new line ending
	sep := []byte(":") // Key: Val

	// for multiline
	ml := []byte("\\\n")  // "\" proceeded with new line
	mlglue := []byte("𐩘") // Unicode character. (U+10A58, &#68184;)

	// Normalize line endings to \n
	buf = bytes.Replace(buf, []byte("\r\n"), nl, -1)

	// Make multiple line params in one line
	// Later we will parse it back to multiline
	buf = bytes.Replace(buf, ml, mlglue, -1)

	// Parse keys: values
	lines := bytes.Split(buf, nl)
	for _, row := range lines {
		row = bytes.TrimSpace(row)

		// Skip not valid format "Key: Val"
		if bytes.Index(row, sep) <= 0 {
			continue
		}

		// Skip comment style rows
		isComment := false ||
			row[0] == "#"[0] || // #
			row[0] == "/"[0] || // // or /*
			row[0] == "-"[0] || // --
			row[0] == "<"[0] || // <!--
			row[0] == "\""[0] || // ""
			row[0] == "~"[0] // ~
		if isComment {
			continue
		}

		// Split to key and val
		prop := bytes.SplitN(row, sep, 2)
		key := bytes.TrimSpace(prop[0])
		val := bytes.TrimSpace(prop[1])

		// Key can't contain spaces
		if bytes.Index(key, []byte(" ")) > 0 {
			continue
		}

		// Is this was multiline make it back to newlines
		// But do not recover "\" at the end
		if bytes.Index(val, mlglue) >= 0 {
			val = bytes.Replace(val, mlglue, nl, -1) // nl NOT ml
		}

		// Assign valid
		params[string(key)] = string(val)
	}

	return params
}

// Parse string (filename) to params
// Example: 1_File name.md
func filenameToParams(fpath string) map[string]string {

	fname := filepath.Base(fpath)
	fname = strings.TrimSpace(fname)

	params := make(map[string]string, 0)
	params["FileName"] = fname
	params["Ext"] = strings.ToLower(filepath.Ext(fname))
	params["IsVisible"] = "Yes"
	// params["Label"] = label - set at the end

	// Is it param file for directory
	dirFname := ".dir"
	if fname == dirFname {
		fname, _ = filepath.Abs(strings.TrimSuffix(fpath, dirFname))
		fname = filepath.Base(fname)
		params["FileName"] = fname

		params["Ext"] = dirFname
		params["IsDir"] = "Yes"
	}

	// Remove extension (can be case sensitive)
	label := strings.TrimSuffix(fname, filepath.Ext(fname)) // note: ext not lowercased

	// SortNr
	// Must be short (because we detect if its not date dd.mm.yyyy)
	arr := strings.SplitN(label, "_", 2)
	if len(arr) >= 2 && len(arr[0]) <= 9 {
		// Check first char for numeric 0-9
		if len(arr[0]) >= 1 && arr[0][0] >= 48 && arr[0][0] <= 57 {
			params["SortNr"] = arr[0]
			label = label[strings.Index(label, "_")+1:] //remove sortNr from label
		}
	}

	// Get limit dates
	// Must be long: dd.mm.yyyy (10)
	// dd.mm.yyyy-dd.mm.yyyy)
	// -dd.mm.yyyy
	// Check first char for numeric 0-9
	arr = strings.SplitN(label, "_", 2)
	if len(arr) >= 2 && len(arr[0]) >= 10 {
		var tFrom, tTo time.Time
		dates := arr[0]
		arr = strings.SplitN(dates, "-", 2)

		var dateErr error

		if len(arr) == 2 {
			// Have both: start and end date
			if len(arr[0]) == 10 {
				// 02.01.2006 => dd.mm.yyyy
				tFrom, dateErr = time.Parse("02.01.2006 15:04", arr[0]+" 00:00")
			}
			if len(arr[1]) == 10 {
				tTo, dateErr = time.Parse("02.01.2006 15:04", arr[1]+" 23:59")
			}
		} else if len(arr) == 1 && len(arr[0]) == 10 {
			// Only one: start date
			tFrom, dateErr = time.Parse("02.01.2006 15:04", arr[0]+" 00:00")
		}

		if dateErr == nil {
			params["DateFrom"] = tFrom.String()
			params["DateTo"] = tTo.String()
			// fmt.Println("-------------", label)
			label = label[strings.Index(label, "_")+1:] //remove date from label

			// Is visible or not
			if tFrom.Year() > 2000 && time.Since(tFrom).Seconds() < 0 {
				params["IsVisible"] = "No"
			}
			if tTo.Year() > 2000 && time.Since(tTo).Seconds() >= 0 {
				params["IsVisible"] = "No"
			}
		}
	}

	// Set before visibility check
	params["Label"] = label
	params["Title"] = label // must be overwritten in fileToParams if have such

	// Slug
	reg, _ := regexp.Compile("[^a-zA-Zā-žĀ-Žа-яА-Я0-9]+") // en, lv, ru
	slug := reg.ReplaceAllString(label, "-")
	slug = ToASCII(slug)
	slug = strings.ToLower(strings.Trim(slug, "-"))
	if slug == "" {
		slug = label
	}
	params["Slug"] = slug

	// Visibility check - always leave it as last check
	// Not visible if no filename
	if params["Label"] == "" {
		params["IsVisible"] = "No"
	}

	// not visible if existing extension not .md
	if params["Ext"] != "" && params["Ext"] != ".md" && params["IsDir"] != "Yes" {
		params["IsVisible"] = "No"
	}

	// Filenames that starts with "." and "~" not visible (also ends with "~")
	if fname[:1] == "." || fname[:1] == "~" || fname[len(fname)-1:] == "~" {
		if params["IsDir"] != "Yes" {
			params["IsVisible"] = "No"
		}
	}

	return params
}
