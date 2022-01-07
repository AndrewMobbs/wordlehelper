package main

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

func (f *freqTable) score(wordlist []string) {
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
	// Then generate word frequency score
	for _, w := range wordlist {
		fScore := 0
		for _, c := range w {
			fScore += distCount[c]
		}
		f.FreqDist[w] = fScore
	}
	f.Total = totalChar
	return
}

func (f *freqTable) findTopScore() string {
	max := 0
	wordChoice := ""
	for w, v := range f.FreqDist {
		if v > max {
			max = v
			wordChoice = w
		}
	}
	return wordChoice
}
