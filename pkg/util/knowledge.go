package util

import (
	"log"
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
			log.Printf("found token: %s", arg)
			return true
		}
	}
	return false
}

func StripPunctuation(s string) string {
	lastChar := s[len(s)-1]
	if unicode.IsPunct(rune(lastChar)) {
		return s[:len(s)-1]
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

		arg = StripPunctuation(arg)
		normalized[strings.ToLower(arg)] = arg
	}
	return normalized
}
