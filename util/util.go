package util

import (
	"fmt"
	"io"
	"os"
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
