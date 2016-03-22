package mango

import gounidecode "github.com/fiam/gounidecode/unidecode" // Rūķis => Rukis

// ToASCII - UTF to Ascii characters
func ToASCII(str string) string {
	return gounidecode.Unidecode(str) // Rūķis => Rukis
}

// Used to merge multiple maps (map[string]string)
// return merged map
func mergeParams(mainMap map[string]string, maps ...map[string]string) map[string]string {
	for _, newMap := range maps {
		for key, val := range newMap {
			_, isKey := mainMap[key]
			if !isKey {
				mainMap[key] = val
			}
		}
	}
	return mainMap
}
