package tiktoken

import (
	"strings"
	"unicode"
)

func CountTokens(s string) int {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	return len(words) + (len(s) / 10)
}
