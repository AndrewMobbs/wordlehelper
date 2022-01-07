package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

func scanFile(file *os.File, wordLength int) []string {
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
		list = append(list, word)
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

func main() {
	filename := os.Args[1]
	file, err := checkOpen(filename)
	if err != nil {
		log.Fatal(err)
	}
	wordlist := scanFile(file, 5)
	f := NewFreqTable()
	f.score(wordlist)
	topWord := f.findTopScore()
	fmt.Printf("Top word is %s\n", topWord)
	return
}
