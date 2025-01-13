package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},

		// custom tests:

		// Emoji 🇬🇧 contains two regional indicator symbols:
		// 1. 🇬 (rune U+1F1EC - regional indicator G)
		// 2. 🇧 (rune U+1F1E7 - regional indicator B)
		{input: "🇬🇧3", expected: "🇬🇧🇬🇧🇬🇧"},

		// emoji: 👨‍👩‍👧‍👦
		// contains 7 runes:
		// 👨 (man)
		// ZWJ (Zero Width Joiner)
		// 👩 (woman)
		// ZWJ
		// 👧 (girl)
		// ZWJ
		// 👦 (boy)
		{input: "👨‍👩‍👧‍👦3", expected: "👨‍👩‍👧‍👦👨‍👩‍👧‍👦👨‍👩‍👧‍👦"},

		{input: "水2😊3한2", expected: "水水😊😊😊한한"},
		{input: "⭐️2abc3💫", expected: "⭐️⭐️abccc💫"},
		{input: "∑2∏3", expected: "∑∑∏∏∏"},

		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
