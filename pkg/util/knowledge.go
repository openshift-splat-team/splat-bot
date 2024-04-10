package util

import (
	"strings"
	"unicode"
)

// TokensPresentAND checks if all of the args are present in the argMap
func TokensPresentAND(argMap map[string]string, args ...string) bool {
	matchedArgs := map[string]bool{}
	for _, arg := range args {
		arg = strings.ToLower(arg)
		if _, exists := argMap[arg]; exists {
			matchedArgs[arg] = true
		}
		if len(matchedArgs) == len(args) {
			return true
		}
	}
	return false
}

// TokensPresentOR checks if any of the args are present in the argMap
func TokensPresentOR(argMap map[string]string, args ...string) bool {
	for _, arg := range args {
		if _, exists := argMap[strings.ToLower(arg)]; exists {
			return true
		}
	}
	return false
}

func StripPunctuation(s string) string {
	lenOnEntry := len(s)
	if len(s) < 2 {
		return s
	}
	firstChar := s[0]
	if unicode.IsPunct(rune(firstChar)) {
		s = s[1:]
	}
	lastChar := s[len(s)-1]
	if unicode.IsPunct(rune(lastChar)) {
		s = s[:len(s)-1]
	}

	if lenOnEntry > len(s) {
		s = StripPunctuation(s)
	}
	return s
}

// NormalizeTokens convert all tokens to lower case for case insensitive matching
func NormalizeTokens(args []string) map[string]string {
	normalized := map[string]string{}
	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}

		arg = strings.ToLower(StripPunctuation(arg))
		arg = strings.TrimSpace(arg)
		normalized[arg] = arg
	}
	return normalized
}

// NormalizeTokensToSlice convert all tokens to lower case for case insensitive matching
func NormalizeTokensToSlice(args []string) []string {
	tokenMap := NormalizeTokens(args)
	normalized := make([]string, len(tokenMap))

	i := 0
	for k := range tokenMap {
		normalized[i] = k
		i++
	}
	return normalized
}
