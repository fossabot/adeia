package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripChars(t *testing.T) {
	testCases := []struct {
		in    string
		chars string
		want  string
	}{
		{"test1234", "+-/=", "test1234"},
		{"  test1234", " ", "test1234"},
		{"  +-/=", "+-/=", "  "},
		{"+-/=", "+-/=", ""},
		{"", "+-/=", ""},
		{"test1234", "test1234", ""},
		{"test1234", "ttt", "es1234"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			got := StripChars(tc.in, tc.chars)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStripWhitespace(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{"test", "test"},
		{"  test1234", "test1234"},
		{"  test1234  ", "test1234"},
		{"\ttest1234\v", "test1234"},
		{"\ttest1234\n", "test1234"},
		{"\ftest1234\r", "test1234"},
		{"\u0085a\u00A0", "a"},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			got := StripWhitespace(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}
