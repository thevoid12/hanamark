package util

import "strings"

func CleanSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
