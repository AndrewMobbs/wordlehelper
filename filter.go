package main

import (
	"fmt"
	"log"
	"strings"
)

// Filter object
// yellow letters are recorded against which round they occur in and position
// green letters are recorded just against position in word
// grey letters are just a flat string
type Filter struct {
	yellow [][wordLength]rune
	green  [wordLength]rune
	grey   string
}

func NewFilter() *Filter {
	filter := Filter{nil, [wordLength]rune{}, ""}
	filter.yellow = make([][5]rune, 1)
	for i := 0; i < wordLength; i++ {
		filter.green[i] = '.'
	}
	return &filter
}

//newYellow adds a new round to the yellow letter filter
func (f *Filter) newYellow() {
	f.yellow = append(f.yellow, [wordLength]rune{})
	k := len(f.yellow)
	for i := 0; i < wordLength; i++ {
		f.yellow[k-1][i] = '.'
	}
}

// processRound takes a single Wordle "round" input of word,result
// and updates the Filter with the consequences of the result
func (f *Filter) processRound(r string) error {
	s := strings.Split(r, ",")
	if len(s) != 2 || len(s[0]) != wordLength || len(s[1]) != wordLength {
		return fmt.Errorf("invalid entry for round %s", r)
	}
	f.newYellow()
	// Validate word
	word := strings.ToLower(s[0])
	if !isAlpha(word) {
		return fmt.Errorf("invalid word in round %s", s[0])
	}
	rs := []rune(word)
	// iterate over result for this word
	for i, v := range s[1] {
		switch v {
		case 'y':
			f.yellow[len(f.yellow)-1][i] = rs[i]
		case 'g':
			if f.green[i] != rs[i] && f.green[i] != '.' {
				return fmt.Errorf("inconsistent green letter specified in round %s", string(v))
			}
			f.green[i] = rs[i]
		case 'x':
			if strings.ContainsRune(string(f.green[:]), rs[i]) {
				// If this is a known green letter marked grey then its not allowed in this position, so treat it as yellow
				f.yellow[len(f.yellow)-1][i] = rs[i]
			} else {
				if !strings.ContainsRune(f.yellowString(), rs[i]) {
					f.grey += string(rs[i])
				}
			}
		default:
			return fmt.Errorf("invalid round result input %s", string(v))
		}
	}
	return nil
}

// getFilter returns a Filter based on a set of rounds of Wordle results
func getFilter(rounds []string) *Filter {
	filter := NewFilter()
	for _, v := range rounds {
		err := filter.processRound(v)
		if err != nil {
			log.Fatal(err)
		}
	}
	return filter
}

// checkWord tests a word against current filter and returns a bool
// indicating whether the word is valid or not, given the filter
func (f *Filter) checkWord(word string) bool {
	//Grey check -- must not have any known grey letters
	if strings.ContainsAny(word, f.grey) {
		return false
	}
	for i, v := range word {
		//Green check -- must contain known green letters in correct place
		if f.green[i] != '.' && f.green[i] != v {
			return false
		}
		//Yellow check -- must contain all yellow letters, but not in yellow positions
		for _, y := range f.yellow[1:] {
			for j, r := range y {
				// Check if it's a known yellow
				if r != 0 && r != '.' {
					// Must contain this rune somewhere
					if !strings.ContainsRune(word, r) {
						return false
					}
					// Must _not_ contain this rune in this position
					if v == r && i == j {
						return false
					}
				}
			}
		}
	}
	return true
}

// yellowString returns the yellow letters for all rounds as a flat string
func (f *Filter) yellowString() string {
	yellowString := ""
	for _, y := range f.yellow {
		for _, r := range y {
			if r != '.' {
				yellowString += string(r)
			}
		}
	}
	return yellowString
}

func (f *Filter) filterList(list []string) []string {
	result := make([]string, 0, len(list))
	for _, word := range list {
		if f.checkWord(word) {
			result = append(result, word)
		}
	}
	return result
}
