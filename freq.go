package main

import (
	"strings"
)

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

type letterFreq map[rune]int

// distCount does a frequency analysis of the letters in the word list
func distCount(wordlist []string) (letterFreq, int) {
	d := make(letterFreq)
	totalChar := 0
	// First do a letter count
	for _, w := range wordlist {
		for _, c := range w {
			totalChar++
			if v, ok := d[c]; ok {
				d[c] = v + 1
			} else {
				d[c] = 1
			}
		}
	}
	return d, totalChar
}

// score takes a list of words and a Filter,
// and applies a score to each word based on the letter frequency and filter
func (f *freqTable) score(wordlist []string, filter *Filter) {
	lf, totalChar := distCount(wordlist)
	green := string(filter.green[:])
	yellow := filter.yellowString()
	// Then generate word frequency score
	for _, w := range wordlist {
		fScore := 0
		seen := ""
		hasRepeat := false
		for _, c := range w {
			if strings.ContainsRune(seen, c) {
				hasRepeat = true
				// Green and Yellow are allowed to repeat, otherwise not. Grey should already have been filtered
				if strings.ContainsRune(yellow, c) || strings.ContainsRune(string(green), c) {
					fScore += lf[c]
				} else {
					// Words with repeated letters that are unknown are not favoured, as we want to explore search space
					fScore = 0
					break
				}
			} else {
				fScore += lf[c]
			}
			seen = seen + string(c)
		}
		f.FreqDist[w] = fScore
		// Hack to strongly prefer no-repeat words over words with repeats
		if !hasRepeat && fScore > 0 {
			f.FreqDist[w] += totalChar
		}
	}
	f.Total = totalChar
}

//holds key/value pairs where key is a string, value an int
type kv struct {
	key   string
	value int
}

// sorted returns the a list of key,value pairs for the  scored words sorted by score, descending
func (f *freqTable) sorted() []kv {
	// Basic insertion sort, better algorithms are available if this is too expensive
	// in practice, it seems fine.
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
