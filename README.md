# Wordle-Solver
Wordle-Solver

Word lists generate by:
```
curl -q https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt | grep -E '^.{6}$' > ./word-listi-ALL.txt

aspell -d en-wo_accents dump master | aspell -l en expand | tr '[:upper:]' '[:lower:]' | grep -E '^[a-z]{5}$' | sort > ./word-list.txt
```

## Usage
```
go run ./wordle-solver.go -help

  -known-letters string
        A string of letters known to be in the word, but their position is unknown (order doesn't matter)

  -known-nonletters string
        A string of letters known to NOT be in the word

  -known-positions string
        A string with the correct letters in their correct positions. Using '_' for unknown positions (default "_____")

  -word-list-path string
        Word List path (default "./word-list.txt")
```
