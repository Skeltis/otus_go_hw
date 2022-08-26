package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	json "github.com/mailru/easyjson"
)

// easyjson:json
type User struct {
	Email []byte `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return getUsersAndCountDomainEntries(r, domain)
}

func getUsersAndCountDomainEntries(reader io.Reader, domainName string) (DomainStat, error) {
	user := &User{}
	result := make(DomainStat)
	var found [][]byte
	var err error

	// usage of methods from strings package actually won't save much more cpu, than precompiled regex on byte stream,
	// but will consume more memory instead
	regexpString := fmt.Sprintf("^(?:[a-zA-Z0-9.!#$%%&'*+\\=?^_`{|}~-]+)"+
		"@([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\\.%s)$", strings.ToLower(domainName))
	compiledRegexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}

	// making rid of using deprecated package, using bufio instead
	scanner := bufio.NewScanner(reader)

	// replacing reading of the whole file with buffered line reader
	for scanner.Scan() {
		// need to convert bytes to lower exactly here, to let correctly pass through regex tree, despite the domain case
		if err = json.Unmarshal(toLowerBytes(scanner.Bytes()), user); err != nil {
			continue
		}

		found = compiledRegexp.FindSubmatch(user.Email)
		if len(found) < 2 {
			continue
		}

		result[string(found[1])]++
	}

	return result, nil
}

// toLowerBytes converts ASCII bytes to lower case.
func toLowerBytes(input []byte) []byte {
	for index, curByte := range input {
		if curByte >= 'A' && curByte <= 'Z' {
			input[index] = curByte + 'a' - 'A'
		}
	}
	return input
}
