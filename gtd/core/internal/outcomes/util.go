package outcomes

import "strings"

func gatherContexts(_description string) (contexts []string) {
	for _, word := range strings.Fields(_description) {
		if strings.HasPrefix(word, "@") && !contains(contexts, word[1:]) {
			contexts = append(contexts, word[1:])
		}
	}
	return contexts
}

func contains(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}
