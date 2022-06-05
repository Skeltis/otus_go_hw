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
	wordFrequenciesArray := packIntoWordFreqsArray(&rawWords)
	sort.Slice(wordFrequenciesArray, func(i, j int) bool {
		if wordFrequenciesArray[i].count > wordFrequenciesArray[j].count {
			return true
		}
		if wordFrequenciesArray[i].count < wordFrequenciesArray[j].count {
			return false
		}
		return wordFrequenciesArray[i].word < wordFrequenciesArray[j].word
	})
	totalRecords := len(wordFrequenciesArray)
	if totalRecords > 10 {
		totalRecords = 10
	}
	top10 := make([]string, totalRecords)
	for index, wordFreqItem := range wordFrequenciesArray[0:totalRecords] {
		top10[index] = wordFreqItem.word
	}
	return top10
}

func packIntoWordFreqsArray(rawWords *[]string) []wordFrequency {
	for _, word := range *rawWords {
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

type wordFrequency struct {
	word  string
	count int
}
