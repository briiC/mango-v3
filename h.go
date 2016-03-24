package mango

import gounidecode "github.com/fiam/gounidecode/unidecode" // Rūķis => Rukis

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
	// m := mainMap - not in new address

	// actual merge
	for _, submap := range maps {
		for key, val := range submap {

			// Set if empty
			if _, isKey := m[key]; !isKey {
				// Not key yet, assign
				m[key] = val
			}

			// Append +Key: Val
			if mval, isKey := m["+"+key]; isKey {
				m[key] += mval
				delete(m, "+"+key) // remove param that appended (+Key)
			}
		}
	}
	return m
}
