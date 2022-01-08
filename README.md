# wordlehelper

Wordlehelper is a helper program for the Wordle word game, and clones. It provides a scored list of suggestions based on letter frequency analysis across a word list, filtered for information from zero or more rounds of the Wordle game.
## Usage
```
Usage: ./wordlehelper <wordlist> [result...]
          <wordlist> is a text file of words to use, one per line
          [result...] is zero or more pairs of word then round result, separated by a comma.
          x = grey (letter not in word)
          y = yellow (letter in word, not in this position)
          g = green (letter in word, in this position)

Example: ./wordlehelper words.txt arose,xxxgy chess,xxggg
```