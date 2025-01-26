package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	punctuationChars  = "-.,;:!?()'"
	onlyPunctuationRe = regexp.MustCompile(`^[` + regexp.QuoteMeta(punctuationChars) + `]+$`)
)

func Top10(in string) []string {
	const topSize = 10
	result := make([]string, 0, topSize)

	wordFrequency := make(map[string]int)
	for _, word := range splitToWords(in) {
		wordFrequency[word]++
	}

	wordsListByFrequency := make(map[int][]string)
	for word, freq := range wordFrequency {
		wordsListByFrequency[freq] = append(wordsListByFrequency[freq], word)
	}

	frequenciesList := make([]int, 0)
	for freq := range wordsListByFrequency {
		frequenciesList = append(frequenciesList, freq)
	}

	sort.Slice(frequenciesList, func(i, j int) bool {
		return frequenciesList[i] > frequenciesList[j]
	})

	for _, freq := range frequenciesList {
		words := wordsListByFrequency[freq]
		sort.Strings(words)

		needToAdd := topSize - len(result)
		if needToAdd == 0 {
			break
		}
		if len(words) < needToAdd {
			result = append(result, words...)
			continue
		}

		result = append(result, words[:needToAdd]...)
		break
	}

	return result
}

// splitToWords split input to slice of words.
// Слово - это либо набор не знаков препинания (все знаки препинания мы удаляем).
// Либо набора знаков препинания, но больше 1.
func splitToWords(in string) []string {
	words := strings.Fields(in)

	i := 0
	for i < len(words) {
		// words with len > 1, that consist only of punctuation marks should be processed as is
		if len(words[i]) > 1 && onlyPunctuationRe.MatchString(words[i]) {
			i++
			continue
		}

		words[i] = strings.ToLower(cutPunctuationMarks(words[i]))
		// skip empty words
		if words[i] == "" {
			words = append(words[:i], words[i+1:]...)
			continue
		}

		i++
	}

	return words
}

func cutPunctuationMarks(in string) string {
	isPunctuation := func() func(r rune) bool {
		firstRune := true

		return func(r rune) bool {
			if !firstRune {
				return false
			}

			firstRune = false
			return strings.ContainsRune(punctuationChars, r)
		}
	}

	leftTrimmed := strings.TrimLeftFunc(in, isPunctuation())
	rightTrimmed := strings.TrimRightFunc(leftTrimmed, isPunctuation())

	if len(rightTrimmed) == 1 && strings.ContainsAny(rightTrimmed, punctuationChars) {
		return ""
	}

	return rightTrimmed
}
