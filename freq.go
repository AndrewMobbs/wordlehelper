package main

import "strings"

type freqTable struct {
	FreqDist map[string]int
	Total    int
}

func NewFreqTable() *freqTable {
	freqScore := make(map[string]int)
	return &freqTable{
		FreqDist: freqScore,
		Total:    0,
	}
}

func (f *freqTable) score(wordlist []string, filter *Filter) {
	distCount := make(map[rune]int)
	totalChar := 0
	// First do a letter count
	for _, w := range wordlist {
		for _, c := range w {
			totalChar++
			if v, ok := distCount[c]; ok {
				distCount[c] = v + 1
			} else {
				distCount[c] = 1
			}
		}
	}
	green := string(filter.green[:])
	yellow := filter.yellowString()
	// Then generate word frequency score
	for _, w := range wordlist {
		fScore := 0
		seen := ""
		for _, c := range w {
			// Green and Yellow are allowed to repeat, otherwise not. Grey should already have been filtered
			if strings.ContainsRune(yellow, c) || strings.ContainsRune(string(green), c) || !strings.ContainsRune(seen, c) {
				fScore += distCount[c]
				seen = seen + string(c)
			} else {
				fScore = 0
				break
			}
		}
		f.FreqDist[w] = fScore
	}
	f.Total = totalChar
}

type kv struct {
	key   string
	value int
}

func (f *freqTable) sorted() []kv {
	// Basic insertion sort, fix if this is too expensive
	sorted := make([]kv, 0, len(f.FreqDist))
	sorted = append(sorted, kv{"", -1})
	for k, v := range f.FreqDist {
		for i, e := range sorted {
			if v > e.value {
				rear := append([]kv{}, sorted[i:]...)
				sorted = append(sorted[0:i], kv{k, v})
				sorted = append(sorted, rear...)
				break
			}
		}
	}
	return sorted
}
