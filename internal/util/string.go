package util

import (
	"strings"
	"unicode"
)

// StripWhitespace strips all Unicode whitespaces from the given string. There are multiple
// methods to achieve this. But according to benchmarks (see: https://stackoverflow.com/a/32081891),
// this is the fastest method.
func StripWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// StripChars strips all runes in `chars` from the given string. For example,
// if `input` is "test1234" and `chars` is "1234", the output will be "test".
func StripChars(input string, chars string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(chars, r) < 0 {
			return r
		}
		return -1
	}

	return strings.Map(filter, input)
}
