package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidString           = errors.New("invalid string")
	isLineStartsWithDigitRegex = regexp.MustCompile(`^\d+`)
	// first part - escaped characters, that are not backslash or digit, like \n,
	// except line start, cause in the middle of the string \\n is correct,
	// second part - escaped character that not backslash or digit in line start
	// third - sequential digits, that aren't escaped digit and digit.
	isIncorrectNumbersExistsRegex = regexp.MustCompile(`[^\\]\\[^\\\d]|^\\[^\\\d]|[^\\]\d{2,}`)
	symbolsSubsequenceRegex       = regexp.MustCompile(`\\\d{1,2}|\D\d|\\{2}\d?`)
)

var inputString string

func Unpack(input string) (string, error) {
	// Place your code here.
	inputString = input

	if inputString == "" {
		return "", nil
	}

	if isLineStartsWithDigitRegex.Match([]byte(inputString)) || isIncorrectNumbersExistsRegex.Match([]byte(inputString)) {
		return "", ErrInvalidString
	}

	symbolsInfo := extractSymbolsMultiplierInfos()
	return convertString(&symbolsInfo), nil
}

func convertString(stringInfo *[]symbolMultiplierInfo) string {
	var strBuilder strings.Builder
	var caretPos int
	for _, symbolInfo := range *stringInfo {
		if symbolInfo.startIndex-caretPos > 0 {
			strBuilder.WriteString(inputString[caretPos:symbolInfo.startIndex])
		}

		strBuilder.WriteString(strings.Repeat(symbolInfo.character, symbolInfo.repeatNumber))
		caretPos = symbolInfo.endIndex
	}
	if caretPos < utf8.RuneCountInString(inputString) {
		strBuilder.WriteString(inputString[caretPos:])
	}
	return strBuilder.String()
}

func extractSymbolsMultiplierInfos() []symbolMultiplierInfo {
	subsequences := symbolsSubsequenceRegex.FindAllSubmatchIndex([]byte(inputString), -1)
	symbInfo := make([]symbolMultiplierInfo, len(subsequences))

	for index, loc := range subsequences {
		symbInfo[index] = extractSymbolInfo(loc[0], loc[1])
	}
	return symbInfo
}

func extractSymbolInfo(start int, end int) symbolMultiplierInfo {
	var repeat int
	subsequence := inputString[start:end]
	symbol, single := extractSymbol(subsequence)
	if single {
		repeat = 1
	} else {
		repeat, _ = strconv.Atoi(subsequence[utf8.RuneCountInString(subsequence)-1:])
	}
	return symbolMultiplierInfo{
		startIndex:   start,
		endIndex:     end,
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
	startIndex   int
	endIndex     int
	repeatNumber int
	character    string
}
