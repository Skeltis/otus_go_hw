package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	counterDictionary map[string]int
	searcher          = regexp.MustCompile(`\s\-\s|[\s.;:,\!\?\\\/()#^%@'"` + "`" + `~+=*]+`)
)

func Top10(input string) []string {
	// Place your code here.
	counterDictionary = make(map[string]int)
	rawWords := searcher.Split(input, -1)
	array := packIntoWordFreqsArray(rawWords)
	sortWordFrequencies(array)
	totalRecords := len(array)
	if totalRecords > 10 {
		totalRecords = 10
	}
	top10 := make([]string, totalRecords)
	for index, pair := range array[0:totalRecords] {
		top10[index] = pair.word
	}
	return top10
}

func packIntoWordFreqsArray(rawWords []string) []wordFrequency {
	for _, word := range rawWords {
		if word != "" {
			lowerWord := strings.ToLower(word)
			counterDictionary[lowerWord]++
		}
	}
	wordFreqs := make([]wordFrequency, len(counterDictionary))
	index := 0
	for key, value := range counterDictionary {
		wordFreqs[index] = wordFrequency{
			word:  key,
			count: value,
		}
		index++
	}
	return wordFreqs
}

func sortWordFrequencies(wordFreqs []wordFrequency) {
	sort.Slice(wordFreqs, func(i, j int) bool {
		if wordFreqs[i].count > wordFreqs[j].count {
			return true
		}
		if wordFreqs[i].count < wordFreqs[j].count {
			return false
		}
		return wordFreqs[i].word < wordFreqs[j].word
	})
}

type wordFrequency struct {
	word  string
	count int
}
