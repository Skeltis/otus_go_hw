package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	// Place your code here.

	if input == "" {
		return "", nil
	}

	if isIncorrectNumbersExists(input) || isLineStartsWithDigit(input) {
		return "", ErrInvalidString
	}

	result := extractSymbolsMultiplierInfos(input)
	return convertString(input, result), nil
}

func isLineStartsWithDigit(line string) bool {
	// line starts with digit
	searcher := regexp.MustCompile(`^\d+`)
	return searcher.Match([]byte(line))
}

func isIncorrectNumbersExists(line string) bool {
	// first part - escaped characters, that are not backslash or digit, like \n,
	// except line start, cause in the middle of the string \\n is correct,
	// second part - escaped character that not backslash or digit in line start
	// third - sequental digits, that aren't escaped digit and digit
	searcher := regexp.MustCompile(`[^\\]\\[^\\\d]|^\\[^\\\d]|[^\\]\d{2,}`)
	return searcher.Match([]byte(line))
}

func convertString(line string, stringInfo []symbolMultiplierInfo) string {
	var strBuilder strings.Builder
	var caretePos uint32
	for _, symbolInfo := range stringInfo {
		if symbolInfo.startIndex-caretePos > 0 {
			strBuilder.WriteString(line[caretePos:symbolInfo.startIndex])
		}

		strBuilder.WriteString(strings.Repeat(symbolInfo.character, symbolInfo.repeatNumber))
		caretePos = symbolInfo.endIndex
	}
	if int(caretePos) < utf8.RuneCountInString(line) {
		strBuilder.WriteString(line[caretePos:])
	}
	return strBuilder.String()
}

func extractSymbolsMultiplierInfos(line string) []symbolMultiplierInfo {
	searcher := regexp.MustCompile(`\\\d{1,2}|\D\d|\\{2}\d?`)
	data := searcher.FindAllSubmatchIndex([]byte(line), -1)
	symbInfo := make([]symbolMultiplierInfo, len(data))

	for index, loc := range data {
		symbInfo[index] = extractSymbolInfo(line, loc[0], loc[1])
	}
	return symbInfo
}

func extractSymbolInfo(line string, start int, end int) symbolMultiplierInfo {
	var repeat int
	subsequence := line[start:end]
	symbol, single := extractSymbol(subsequence)
	if single {
		repeat = 1
	} else {
		repeat, _ = strconv.Atoi(subsequence[utf8.RuneCountInString(subsequence)-1:])
	}
	return symbolMultiplierInfo{
		startIndex:   uint32(start),
		endIndex:     uint32(end),
		repeatNumber: repeat,
		character:    symbol,
	}
}

func extractSymbol(sequence string) (string, bool) {
	if strings.HasPrefix(sequence, `\`) {
		if utf8.RuneCountInString(sequence) == 2 {
			return sequence[1:], true
		}
		return sequence[1 : utf8.RuneCountInString(sequence)-1], false
	}
	return sequence[0 : utf8.RuneCountInString(sequence)-1], false
}

type symbolMultiplierInfo struct {
	startIndex   uint32
	endIndex     uint32
	repeatNumber int
	character    string
}
