package util

import "strings"

func LowerEqual(src string, dest string) bool {
	if len(src) != len(dest) {
		return false
	}
	return strings.ToLower(src) == dest
}

func UpperEqual(src string, dest string) bool {
	if len(src) != len(dest) {
		return false
	}
	return strings.ToUpper(src) == dest
}
