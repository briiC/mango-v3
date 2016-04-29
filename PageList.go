package mango

import (
	"log"
	"math"
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
	iNum, _ := strconv.Atoi(pages[i].Get("SortNr"))
	jNum, _ := strconv.Atoi(pages[j].Get("SortNr"))

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
		p.Set("SortNr", strconv.Itoa(rand.Intn(count*10)+1))
	}

	// Sort now by default
	sort.Sort(pages)
}

// Sort list by given sortType
func (pages PageList) Sort(sortType string) {
	if len(pages) >= 2 {
		switch sortType {
		case "Reverse":
			sort.Sort(sort.Reverse(pages))
		case "Random":
			pages.Randomize()
		default:
			sort.Sort(pages)
		}
	}
}

// Paging - Make paging params
// Returns necessary params for paging
func (pages PageList) Paging(pNum, pSize, pLimit int) (PageList, map[string]int) {
	len := pages.Len()
	pageCount := int(math.Ceil(float64(len) / float64(pSize)))

	if pSize < 1 {
		pSize = 1
	}

	if pLimit > 0 && pageCount > pLimit {
		pageCount = pLimit
		len = pSize * pageCount
	} else {
		pLimit = -1
	}

	if pNum > pageCount {
		pNum = pageCount
	}
	if pNum < 1 {
		pNum = 1
	}

	pPrev := pNum - 1
	pNext := pNum + 1
	if pNext > pageCount {
		pNext = 0
	}

	from := (pNum - 1) * pSize
	to := pNum * pSize

	if from >= len {
		from = len - pSize
	}
	if to > len {
		to = len
	}

	if len <= 0 {
		from = 0
		to = 0
	}

	pages = pages[from:to]

	return pages, map[string]int{
		"PPrev":       pPrev,
		"PNum":        pNum,
		"PNext":       pNext,
		"PSize":       pSize,
		"PFrom":       from,
		"PTo":         to,
		"PTotalPages": int(math.Ceil(float64(len) / float64(pSize))),
		"PTotalItems": len,
		"PLimit":      pLimit,
	}
}

// Print pages in list
func (pages PageList) Print() {
	log.Printf("--- %d pages ------------------------------------------------", len(pages))
	for _, p := range pages {
		// log.Printf("- %20s", p.Get("Slug"))
		p.PrintRow()
	}
	log.Println("---------------------------------------------------")
}
