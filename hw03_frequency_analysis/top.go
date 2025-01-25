package hw03frequencyanalysis

import (
	"sort"
	"strings"
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

func splitToWords(in string) []string {
	words := strings.Fields(in)

	i := 0
	for i < len(words) {
		words[i] = strings.ToLower(strings.Trim(words[i], ".,;:!?-()'\""))
		if words[i] == "" {
			words = append(words[:i], words[i+1:]...)
			continue
		}

		i++
	}

	return words
}
