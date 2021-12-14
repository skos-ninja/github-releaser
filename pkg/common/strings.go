package common

import "strings"

// Contains iterates over array of strings and returns true if any of them match (case insensitive)
func Contains(arr []string, str string) bool {
	for _, s := range arr {
		if strings.Contains(strings.ToLower(s), strings.ToLower(str)) {
			return true
		}
	}

	return false
}
