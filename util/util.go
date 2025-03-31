package util

import (
	"path/filepath"
	"strings"
)

func CleanSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func RemoveExtentionFromFile(path string) string {
	ext := filepath.Ext(path)
	path = path[0 : len(path)-(len(ext))]
	return path
}
