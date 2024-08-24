package tiktoken

import (
	"strings"
	"unicode"
)

// CountTokens estimates the number of tokens in a string.
// This is a simplified version and may not be as accurate as GPT's tiktoken.
func CountTokens(s string) int {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Estimate: 1 token per word, plus some extra for punctuation and whitespace
	return len(words) + (len(s) / 10)
}
