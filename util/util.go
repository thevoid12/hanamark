package util

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

// relURL finds the relative url between a sounce folder and destination file
func RelURL(fromFile, toFile string) (string, error) {
	// if it is a directory then we are already good to go
	info, err := os.Stat(fromFile)
	var fromDir string
	if err == nil {
		if info.IsDir() {
			fromDir = fromFile
		} else {
			fromDir = filepath.Dir(fromFile)
		}
	} else if os.IsNotExist(err) {
		// If file doesn't exist, assume it's a file path and get its directory
		// This is necessary for generating links from pages that haven't been created yet
		fromDir = filepath.Dir(fromFile)
	} else {
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
	_, err = f.Write([]byte(content))
	if err != nil {
		l.Sugar().Error("writing into the file failed", err)
		return err
	}

	return nil
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
