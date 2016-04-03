package mango

import (
	"errors"
	"log"
	"strings"
	"time"

	gounidecode "github.com/fiam/gounidecode/unidecode" // Rūķis => Rukis
)

// ToASCII - UTF to Ascii characters
func ToASCII(str string) string {
	return gounidecode.Unidecode(str) // Rūķis => Rukis
}

// Used to merge multiple maps (map[string]string)
// return merged map
// NOT THREAD-SAFE. Use this testing heavily.
// Safe if using in application init phase.
func mergeParams(mainMap map[string]string, maps ...map[string]string) map[string]string {
	//copy mainMap for concurrent write
	m := make(map[string]string, 0)
	for key, val := range mainMap {
		m[key] = val
	}

	// actual merge
	for _, submap := range maps {
		for key, val := range submap {

			// Create if empty
			if _, isKey := m[key]; !isKey {
				m[key] = val
			}

			// Append +Key: Val
			if mval, isKey := m["+"+key]; isKey {
				m[key] += ", " + mval
				delete(m, "+"+key) // remove param that appended (+Key)
			}
		}
	}
	return m
}

// Parse any custom string
// 2006-01-02 15:04:05
// 01 - month
// 02 - day
func toTime(s string) (time.Time, error) {
	sLen := len(s)
	dtNow := time.Now()

	// Sep counts
	timeSeps := strings.Count(s, ":")
	dotCount := strings.Count(s, ".")
	dashCount := strings.Count(s, "-")
	slashCount := strings.Count(s, "/")

	var dt time.Time
	var err error

	switch {

	// Only time
	case sLen == 5 && timeSeps == 1:
		dt, err = time.Parse("15:04", s)
	case sLen == 8 && timeSeps == 2:
		dt, err = time.Parse("15:04:05", s)

		// d.m.y
	case dotCount == 2 && timeSeps == 1:
		dt, err = time.Parse("02.01.2006 15:04", s)
	case dotCount == 2 && timeSeps == 2:
		dt, err = time.Parse("02.01.2006 15:04:05", s)
	case dotCount == 1 && timeSeps == 1:
		dt, err = time.Parse("02.01 15:04", s)
	case dotCount == 2 && timeSeps == 0:
		dt, err = time.Parse("02.01.2006", s)
	case dotCount == 1 && timeSeps == 0:
		dt, err = time.Parse("02.01", s)

		// y-m-d
	case dashCount == 2 && timeSeps == 1:
		dt, err = time.Parse("2006-01-02 15:04", s)
	case dashCount == 2 && timeSeps == 2:
		dt, err = time.Parse("2006-01-02 15:04:05", s)
	case dashCount == 1 && timeSeps == 1:
		dt, err = time.Parse("01-02 15:04", s)
	case dashCount == 2 && timeSeps == 0:
		dt, err = time.Parse("2006-01-02", s)
	case dashCount == 1 && timeSeps == 0:
		dt, err = time.Parse("01-02", s)

		// mm/dd/yyyy
	case slashCount == 2 && timeSeps == 1:
		dt, err = time.Parse("01/02/2006 15:04", s)
	case slashCount == 2 && timeSeps == 2:
		dt, err = time.Parse("01/02/2006 15:04:05", s)
	case slashCount == 1 && timeSeps == 1:
		dt, err = time.Parse("01/02 15:04", s)
	case slashCount == 2 && timeSeps == 0:
		dt, err = time.Parse("01/02/2006", s)
	case slashCount == 1 && timeSeps == 0:
		dt, err = time.Parse("01/02", s)

	default:
		err = errors.New("ERROR: \"" + s + "\" not in correct datetune format")
	}

	if err != nil {
		// On error return what you got
		return dtNow, err
	}

	// Fill unset parts of datetime
	if dt.Year() == 0 {
		if (dotCount+dashCount+slashCount == 0) && timeSeps > 0 {
			// Set date for "only-time" value
			// 0000-01-01
			dt = dt.AddDate(dtNow.Year(), int(dtNow.Month())-1, dtNow.Day()-1)
		} else {
			// Set date for "date-time" value
			// OOnly year is empty
			dt = dt.AddDate(dtNow.Year(), 0, 0)
		}
	}

	return dt, err
}

// Print map in human readable format
func printMap(fname string, m map[string]string) {
	log.Println("---", fname, "--------------------------------------------")
	for key, val := range m {
		log.Printf("%20s: %s \n", key, val)
	}
	log.Println("--------------------------------------------------------")
}
