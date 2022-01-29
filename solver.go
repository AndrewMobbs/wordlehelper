package main

import (
	"math"
	"strings"
)

func getSolution(guessList []string, solList []string, filter *Filter) string {
	guessCount := make(map[string]int)

	for _, solWord := range solList {

		// Algorithm 1, per-round minmax, find guess that excludes most words given current knowledge
		// Need starting filter! Add each round filter to starting filter.
		//
		// For each word on the solution list:
		// For each word on the guess list:
		// run a game, get result, apply filter to solution list, record length of list for guess word (just add to total?)
		// Sort record of games
		// Give guess with overall shortest list lengths.

		for _, guess := range guessList {
			thisFilter := NewFilter()
			thisFilter.green = filter.green
			copy(thisFilter.yellow, filter.yellow)
			thisFilter.grey = filter.grey

			result := getRoundResult(guess, solWord)
			roundResult := guess + "," + result
			thisFilter.processRound(roundResult)
			thisCount := 0
			for _, w := range solList {
				if thisFilter.checkWord(w) {
					thisCount++
				}
			}
			// Is this needed?
			if thisCount == 0 {
				thisCount = len(solList)
			}
			guessCount[guess] += thisCount
		}
	}

	min := math.MaxInt32
	choice := ""
	if len(solList) <= 2 {
		choice = solList[0]
	} else {
		for w, v := range guessCount {
			if v < min && v > 0 {
				choice = w
				min = v
			}
		}
	}
	return choice
}

func getRoundResult(guess string, solution string) string {
	s := []rune(solution)
	r := ""
	for i, c := range guess {
		if c == s[i] {
			r = r + "g"
		} else if strings.ContainsRune(solution, c) {
			r = r + "y"
		} else {
			r = r + "x"
		}
	}
	return r
}
