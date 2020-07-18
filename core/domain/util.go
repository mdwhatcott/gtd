package domain

import "strings"

func gatherContexts(description string) (contexts_ []string) {
	for _, WORD := range strings.Fields(description) {
		if strings.HasPrefix(WORD, "@") && !contains(contexts_, WORD[1:]) {
			contexts_ = append(contexts_, WORD[1:])
		}
	}
	return contexts_
}

func contains(haystack []string, needle string) bool {
	for _, STRAW := range haystack {
		if STRAW == needle {
			return true
		}
	}
	return false
}
