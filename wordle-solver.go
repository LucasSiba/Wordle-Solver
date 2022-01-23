package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"
)

const WORD_LENGTH = 5

var (
	wordListPath    = flag.String("word-list-path", "./word-list.txt", "Word List path")
	knownPositions  = flag.String("known-positions", "_____", "A string with the correct letters in their correct positions. Using '_' for unknown positions")
	knownLetters    = flag.String("known-letters", "", "A string of letters known to be in the word, but their position is unknown (order doesn't matter)")
	knownNonLetters = flag.String("known-nonletters", "", "A string of letters known to NOT be in the word")
)

func init() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	if len(*knownPositions) != WORD_LENGTH {
		log.Fatalf("-known-positions must be exactly %d characters", WORD_LENGTH)
	}

	pass, _ := regexp.MatchString(`^[a-z_]+$`, *knownPositions)
	if pass == false {
		log.Fatalf("-known-positions contains illegal letters")
	}

	pass, _ = regexp.MatchString(`^[a-z]+$`, *knownLetters)
	if *knownLetters != "" && pass == false {
		log.Fatalf("-known-letters contains illegal letters")
	}

	pass, _ = regexp.MatchString(`^[a-z]+$`, *knownNonLetters)
	if *knownNonLetters != "" && pass == false {
		log.Fatalf("-known-nonletters contains illegal letters")
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	wordList, err := readLines(*wordListPath)
	if err != nil {
		log.Fatalf("Failed to read word list: '%s': ", err)
	}

	if *knownPositions == "_____" && *knownLetters == "" && *knownNonLetters == "" {
		printBestGuess(wordList)
		return
	}

	if *knownPositions != "_____" {
		for i, posC := range []byte(*knownPositions) {
			var reducedWordList []string
			if posC == '_' {
				continue
			}
			for _, word := range wordList {
				if word[i] != posC {
					//log.Printf("Skipping '%s' for wrong position\n", word)
					continue
				} else {
					reducedWordList = append(reducedWordList, word)
					//log.Printf("Keeping '%s'\n", word)
				}
			}
			wordList = reducedWordList
		}
	}

	if *knownLetters != "" {
		for _, c := range *knownLetters {
			var reducedWordList []string
			for _, word := range wordList {
				var foundLetter bool
				for _, c2 := range word {
					if c == c2 {
						foundLetter = true
					}
				}
				if !foundLetter {
					//log.Printf("Skipping '%s' for missing letter\n", word)
					continue
				} else {
					reducedWordList = append(reducedWordList, word)
					//log.Printf("Keeping '%s'\n", word)
				}
			}
			wordList = reducedWordList
		}
	}

	if *knownNonLetters != "" {
		for _, c := range *knownNonLetters {
			var reducedWordList []string
			for _, word := range wordList {
				var foundLetter bool
				for _, c2 := range word {
					if c == c2 {
						foundLetter = true
					}
				}
				if foundLetter {
					//log.Printf("Skipping '%s' for found letter\n", word)
					continue
				} else {
					reducedWordList = append(reducedWordList, word)
					//log.Printf("Keeping '%s'\n", word)
				}
			}
			wordList = reducedWordList
		}
	}

	printBestGuess(wordList)
}

func printBestGuess(wordList []string) {
	log.Printf("There are %d remaining words in the word list", len(wordList))

	// Find words without duplicate letters, to get more letter coverage on the guess
	var reducedWordList []string
	for _, word := range wordList {
		// For small WORD_LENGTH, it's faster to compare then make a map
		skipWord := false
		for i, c := range word {
			if i == WORD_LENGTH {
				continue
			}
			for _, c2 := range word[i+1:] {
				if c == c2 {
					skipWord = true
					// log.Printf("Skipping '%s' for duplicates\n", word)
					break
				}
			}
		}
		if !skipWord {
			// log.Printf("Adding '%s' to guesses\n", word)
			reducedWordList = append(reducedWordList, word)
		}
	}

	if len(reducedWordList) < 10 {
		log.Printf("This is all that's left: %+v", wordList)
	}

	if len(reducedWordList) != 0 {
		log.Printf("Here are a couple of good next guesses: %s %s %s",
			reducedWordList[rand.Intn(len(reducedWordList))],
			reducedWordList[rand.Intn(len(reducedWordList))],
			reducedWordList[rand.Intn(len(reducedWordList))])
	}
}
