package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

func isAlpha(word string) bool {
	for _, c := range word {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

const defaultListLength = 100

func scanFile(file *os.File, filter *Filter, wordLength int) []string {
	scanner := bufio.NewScanner(file)
	list := make([]string, 0, defaultListLength)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) != wordLength {
			continue
		}
		word = strings.ToLower(word)
		if !isAlpha(word) {
			continue
		}
		if filter.checkWord(word) {
			list = append(list, word)
		}
	}
	return list
}

func checkOpen(filename string) (*os.File, error) {
	log.Trace("checkOpen()")
	filestat, err := os.Stat(filename)
	exists := true
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			exists = false
		} else {
			log.Fatal("Error statting wordlist file: ", err)
		}
	}

	if !filestat.Mode().IsRegular() {
		log.Fatal("Wordlist is not a regular file")
	}

	if exists {
		return os.Open(filename)
	}
	return nil, err
}

const wordLength = 5
const topResults = 10

func printHelp(name string) {
	fmt.Printf("Usage: %s <wordlist> [result...]\n", name)
	fmt.Printf("          <wordlist> is a text file of words to use, one per line\n")
	fmt.Printf("          [result...] is zero or more pairs of word then round result, separated by a comma.\n")
	fmt.Printf("              x = grey (letter not in word)\n")
	fmt.Printf("              y = yellow (letter in word, not in this position)\n")
	fmt.Printf("              g = green (letter in word, in this position)\n")
	fmt.Printf("\nExample: %s words.txt arose,xxxgy chess,xxggg\n", name)
}

func doWordleHelper(args []string) {
	filename := args[1]
	if filename == "-h" || filename == "--help" {
		printHelp(args[0])
		return

	}
	file, err := checkOpen(filename)
	if err != nil {
		log.Fatal(err)
	}
	filter := getFilter(args[2:])
	wordlist := scanFile(file, filter, wordLength)
	f := NewFreqTable()
	f.score(wordlist, filter)
	for i, e := range f.sorted() {
		if e.key != "" {
			fmt.Printf("#%d word is %s with score %d\n", i+1, e.key, e.value)
		}
		if i == topResults-1 {
			break
		}
	}
}

func doWordleSolver(args []string) string {
	solListFile, err := checkOpen(args[1])
	if err != nil {
		log.Fatal(err)
	}
	guessWordList, err := checkOpen(args[2])
	if err != nil {
		log.Fatal(err)
	}
	filter := getFilter(args[3:])
	// Get all words on the guess list
	guessList := scanFile(guessWordList, NewFilter(), wordLength)
	// Filter the solution list
	solList := scanFile(solListFile, filter, wordLength)
	return getSolution(guessList, solList, filter)

}

func doSolverTest(args []string) {
	solListFile, err := checkOpen(args[1])
	if err != nil {
		log.Fatal(err)
	}
	guessWordList, err := checkOpen(args[2])
	if err != nil {
		log.Fatal(err)
	}
	guessList := scanFile(guessWordList, NewFilter(), wordLength)
	filter := NewFilter()
	solList := scanFile(solListFile, filter, wordLength)
	firstRoundCache := make(map[string]string)
	for _, word := range solList {
		roundSolList := solList
		guess := args[3]
		wordResult := word + ":"
		result := ""
		roundNum := 0
		for len(roundSolList) > 1 && result != "ggggg" {
			roundNum++
			thisFilter := NewFilter()
			thisFilter.green = filter.green
			copy(thisFilter.yellow, filter.yellow)
			thisFilter.grey = filter.grey

			result = getRoundResult(guess, word)
			roundResult := guess + "," + result
			if guess != word {
				wordResult += roundResult + ":"
			}
			thisFilter.processRound(roundResult)
			roundSolList = thisFilter.filterList(roundSolList)
			if v, ok := firstRoundCache[roundResult]; ok && roundNum == 1 {
				guess = v
			} else {
				guess = getSolution(guessList, roundSolList, thisFilter)
				firstRoundCache[roundResult] = guess
			}
		}
		wordResult += guess + ",ggggg"
		fmt.Printf("%s\n", wordResult)
	}

}

func main() {
	if len(os.Args) == 1 {
		printHelp(os.Args[0])
		return
	}
	switch filepath.Base(os.Args[0]) {
	case "wordlehelper":
		doWordleHelper(os.Args)
	case "wordlesolver":
		fmt.Printf("%s\n", doWordleSolver(os.Args))
	case "solvertest":
		if len(os.Args) != 4 {
			fmt.Printf("Usage solvertest <solution list> <word list> <start word>")
			os.Exit(1)
		}
		doSolverTest(os.Args)
	default:
		log.Fatal("Unknown program name: %s", os.Args[0])
	}
}
