package utils

import (
	"fmt"
	"strings"
)

var keywords []string = []string{"add", "character"}

func createEditArray(word string) []string {
	word = strings.ToLower(word)
	characters := "abcdefghijklmnopqrstuvwxyz"
	arr := make([]string, 0)

	for i, _ := range word {
		arr = append(arr, word[:i]+""+word[i+1:]) // delete
		swap := []rune(word)
		if i+1 < len(word) {
			swap[i], swap[i+1] = swap[i+1], swap[i]
			arr = append(arr, string(swap)) // swap
		}
		for _, a := range characters {
			arr = append(arr, word[:i]+string(a)+word[i+1:]) // replace
			arr = append(arr, word[:i]+string(a)+word[i:])   // insert
		}
	}

	return arr
}

func inKeywords(word string) bool {
	for _, keyword := range keywords {
		if word == keyword {
			return true
		}
	}
	return false
}

func Validate(word string, times int) {
	nEdits := createEditArray(word)

	i := 1
	for i < times {
		edits := make([]string, 0)
		for _, w := range nEdits {
			edits = append(edits, createEditArray(w)...)
		}
		nEdits = append(nEdits, edits...)
		i += 1
	}

	for _, possibleCorrection := range nEdits {
		if inKeywords(possibleCorrection) {
			fmt.Printf("Did you mean by %s?\n", possibleCorrection)
			return
		}
	}
}
