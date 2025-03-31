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

func RemoveRootPartOfDir(oldpath, destMDRoot string) string {
	// Normalize destMDRoot to match the format of originalPath
	normalizedRoot := strings.TrimPrefix(destMDRoot, "./")

	res := filepath.Join(".", strings.TrimPrefix(oldpath, normalizedRoot))

	return res
}
