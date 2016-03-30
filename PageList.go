package mango

import "strconv"

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
		return false
	}

	return iNum < jNum
}
