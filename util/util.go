package util

import (
	"context"
	"embed"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

func CleanSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func RemoveExtentionFromFile(path string) string {
	ext := filepath.Ext(path)
	path = path[0 : len(path)-(len(ext))]
	return path
}

// CleanURLPath converts a generated ".../name.html" destination path into the
// extension-less form static hosts with "clean URL" rewriting (e.g. Cloudflare Pages)
// actually serve without a redirect: ".../index.html" -> ".../", ".../name.html" -> ".../name".
func CleanURLPath(p string) string {
	p = filepath.ToSlash(p)
	dir, file := path.Split(p)
	if file == "index.html" {
		return dir
	}
	if strings.HasSuffix(file, ".html") {
		return dir + strings.TrimSuffix(file, ".html")
	}
	return p
}

func RemoveRootPartOfDir(oldpath, destHtmlDir string) string {
	// Normalize destHtmlDir to match the format of originalPath
	normalizedRoot := strings.TrimPrefix(destHtmlDir, "./")
	res := filepath.Join(".", strings.TrimPrefix(oldpath, normalizedRoot))

	return res
}

// relURL finds the relative url between a sounce folder and destination file
func RelURL(fromFile, toFile string) (string, error) {
	// if it is a directory then we are already good to go
	info, err := os.Stat(fromFile)
	var fromDir string
	switch {
	case err == nil:
		if info.IsDir() {
			fromDir = fromFile
		} else {
			fromDir = filepath.Dir(fromFile)
		}
	case os.IsNotExist(err):
		// If file doesn't exist, assume it's a file path and get its directory
		// This is necessary for generating links from pages that haven't been created yet
		fromDir = filepath.Dir(fromFile)
	default:
		return "", err // permission error or other
	}

	rel, err := filepath.Rel(fromDir, toFile)
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(rel), nil
}

// CopyAssets copies images from sourceDir to destDir, preserving the directory structure
func CopyAssets(sourceDir, destDir string) error {
	if sourceDir == "" {
		return nil // skip if no source assets
	}
	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return nil // skip if source does not exist
	}
	// Ensure the destination directory exists
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, they will be created when copying files
		if info.IsDir() {
			return nil
		}

		// Get the relative path from sourceDir
		relPath, err := filepath.Rel(sourceDir, srcPath)
		if err != nil {
			return err
		}

		// Construct the destination path
		destPath := filepath.Join(destDir, relPath)

		// Ensure the parent directory exists in the destination
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create destination subdirectory: %w", err)
		}

		// Copy the image file
		return copyFile(srcPath, destPath)
	})

	if err != nil {
		return fmt.Errorf("error copying images: %w", err)
	}
	return nil
}

func WriteIntoFile(ctx context.Context, content string, filePath string) error {
	l := logs.GetLoggerctx(ctx)

	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		l.Sugar().Error("failed to create directories", err)
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}

	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		l.Sugar().Error("writing into the file failed", err)
		return err
	}

	return nil
}

// CopyFile copies a file from src to dst, replacing if it exists (exported version)
func CopyFile(src, dst string) error {
	return copyFile(src, dst)
}

// copyFile copies a file from src to dst, replacing if it exists
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func ParseTimeFlexible(s string) (time.Time, error) {
	for _, layout := range model.TimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}

func JoinURL(baseUrl, p string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}

func EnsureEmptyDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("directory not empty: %s", path)
	}
	return os.MkdirAll(path, 0755)
}

func CopyEmbedDir(efs embed.FS, src, dst string) error {
	return fs.WalkDir(efs, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == "embed.go" {
			return nil
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		return copyEmbedFile(efs, path, target)
	})
}

func copyEmbedFile(efs embed.FS, src, dst string) error {
	in, err := efs.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err // permission / IO error
	}
	return fi.IsDir(), nil
}

// FindTemplateUpward searches for a template file by traversing upward through directories
// starting from startDir until it finds the file or reaches the rootPath boundary.
// rootPath is the root template directory boundary (e.g., "./configurables/templates")
// fileName is the name of the template file to search for (e.g., "single.html", "list.html")
// Returns the full path to the template file if found, or an error if not found.
func FindTemplateUpward(startDir, rootPath, fileName string) (string, error) {
	// Clean the paths to ensure consistent handling
	currentDir := filepath.Clean(startDir)
	rootPath = filepath.Clean(rootPath)

	// Ensure startDir is within or equal to rootPath
	relPath, err := filepath.Rel(rootPath, currentDir)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("start directory '%s' is not within root path '%s'", startDir, rootPath)
	}

	for {
		// Construct the potential template path
		templatePath := filepath.Join(currentDir, fileName)

		// Check if the file exists and is not a directory
		info, statErr := os.Stat(templatePath)
		if statErr == nil && !info.IsDir() {
			// File found!
			return templatePath, nil
		}

		// Check if we've reached the root boundary
		if currentDir == rootPath {
			// We've reached the root template path without finding the file
			return "", fmt.Errorf("template file '%s' not found in '%s' or any parent directories up to root '%s'", fileName, startDir, rootPath)
		}

		// Move to parent directory
		parentDir := filepath.Dir(currentDir)

		// Safety check: ensure we don't go beyond root (shouldn't happen with above check)
		if parentDir == currentDir {
			return "", fmt.Errorf("reached filesystem root while searching for '%s'", fileName)
		}

		currentDir = parentDir
	}
}

// see if parent directory is actually the parent dir of the child directory
// child dir: blogs/jan/blog1.html
// parent dir: blogs/
// isDirUnder is true
func IsDirUnder(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	if parent == "." {
		return true
	}
	if child == parent {
		return true
	}

	sep := string(os.PathSeparator)
	return strings.HasPrefix(child, parent+sep)
}

// CountWords counts the number of words in a text string
func CountWords(text string) int {
	words := 0
	inWord := false

	for _, char := range text {
		if unicode.IsSpace(char) {
			inWord = false
		} else if !inWord {
			inWord = true
			words++
		}
	}

	return words
}

// CalculateReadTime calculates reading time in minutes based on word count
// Uses average reading speed of 200 words per minute
func CalculateReadTime(text string) int {
	wordCount := CountWords(text)
	readTime := wordCount / 200

	// Minimum 1 minute read time for any content
	if readTime < 1 && wordCount > 0 {
		return 1
	}

	return readTime
}
