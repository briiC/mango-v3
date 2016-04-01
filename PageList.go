package mango

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// PageList is slice as []*Page
type PageList []*Page

// Len is part of sort.Interface.
func (pages PageList) Len() int {
	return len(pages)
}

// Swap is part of sort.Interface.
func (pages PageList) Swap(i, j int) {
	pages[i], pages[j] = pages[j], pages[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (pages PageList) Less(i, j int) bool {
	iNum, _ := strconv.Atoi(pages[i].Params["SortNr"])
	jNum, _ := strconv.Atoi(pages[j].Params["SortNr"])

	if iNum == 0 || jNum == 0 {
		// unset or broken SortNr. Do nothing
		return false
	}

	return iNum < jNum
}

// Randomize slice
// Randomizes param SortNr for all pages and sorts
func (pages PageList) Randomize() {
	count := len(pages)
	rand.Seed(time.Now().UnixNano() + int64(count))

	// Make SortNr as random
	for _, p := range pages {
		p.Params["SortNr"] = strconv.Itoa(rand.Intn(count*10) + 1)
	}

	// Sort now by default
	sort.Sort(pages)
}

// Print pages in list
func (pages PageList) Print() {
	log.Println("---------------------------------------------------")
	for _, p := range pages {
		log.Printf("- %20s", p.Get("Slug"))
	}
	log.Println("---------------------------------------------------")
}
