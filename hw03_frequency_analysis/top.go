package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	Word string
	Freq int
}

func Top10(text string) []string {
	s := strings.Fields(text)
	uniqueWords := make(map[string]int)

	for _, v := range s {
		uniqueWords[v]++
	}

	var words []Word
	for k, v := range uniqueWords {
		words = append(words, Word{k, v})
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].Freq == words[j].Freq {
			return words[i].Word < words[j].Word
		}
		return words[i].Freq > words[j].Freq
	})

	result := make([]string, 0, 10)

	for _, v := range words {
		if len(result) == 10 {
			break
		}
		result = append(result, v.Word)
	}

	return result
}
