# wordlehelper

Wordlehelper is a helper program for the Wordle word game, and clones. It provides a scored list of suggestions based on letter frequency analysis across a word list, filtered for information from zero or more rounds of the Wordle game.

It is best used as a helper rather than a solver. At the moment it tries to suggest the most likely answer given known constraints, not the query that best reduces the search space. 

The first two suggestions from the provided wordlist are "arose" and then "unity" (these can be found by running `wordlehelper words.txt` and `wordlehelper words.txt arose,xxxxx` to pretend that "arose" found no results). After these two guesses, the search space for a given Wordle should be usefully constrained, if not then a useful third choice is "mulch". Maybe a future version will try to suggest search-space pruning guesses as well as solution guesses.

The provided word list isn't guaranteed to be the same as the one used by Wordle (or any given clone). You're free to create any alternative word list you like.

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
