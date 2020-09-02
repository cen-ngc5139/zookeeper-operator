package utils

import "strings"

func Joins(args ...string) string {
	var str strings.Builder
	for _, arg := range args {
		str.WriteString(arg)
	}
	return str.String()
}

// StringsInSlice returns true if the given strings are found in the provided slice, else returns false
func StringsInSlice(strings []string, slice []string) bool {
	asMap := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		asMap[s] = struct{}{}
	}
	for _, s := range strings {
		if _, exists := asMap[s]; !exists {
			return false
		}
	}
	return true
}

// StringInSlice returns true if the given string is found in the provided slice, else returns false
func StringInSlice(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// RemoveStringInSlice returns a new slice with all occurrences of s removed,
// keeping the given slice unmodified
func RemoveStringInSlice(s string, slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return result
}
