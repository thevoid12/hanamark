package parser

import (
	"context"
	"errors"
	logs "hanamark/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// root files need to be updated in the end after all parsing in dest folder is done
func CreateUpdateBasefile(ctx context.Context) error {

	l := logs.GetLoggerctx(ctx)

	rootSrcDir := viper.GetString("filepath.sourceMDRoot")
	rootDestDir := viper.GetString("filepath.destMDRoot")
	// baseFileMap has the list of base files which has subfolder in them so these basefiles
	// should have links to the subfiles in it eg blogs.html should have links of all the blogs
	// in its corresponding subfolder
	baseFileMap := viper.GetStringMapString("fileMeta.baseFilesMap")

	for basefileName, bfdir := range baseFileMap {

		_, err := os.Stat(basefileName)
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Create(bfdir)
			if err != nil {
				l.Sugar().Error("create file failed", err)
				return err
			}
		}

		// traverse through the sub directory of src  and create links to the base file in destination
		err = filepath.Walk(filepath.Join(rootDestDir, bfdir), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Process only Markdown files
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				// Determine relative path from source root
				relPath, err := filepath.Rel(rootSrcDir, path)
				if err != nil {
					return err
				}

				// Construct the corresponding destination path
				destPath := filepath.Join(rootDestDir, relPath)
				destDir := filepath.Dir(destPath)

				// Ensure the destination directory exists
				err = os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					l.Sugar().Error("make destination director failed", err)
					return err
				}

				// Generate markdown with file links
				err = ParseMarkdownToHtml(filepath.Join(rootSrcDir, bfdir), destPath)
				if err != nil {
					l.Sugar().Error("Error parsing markdown to html", err)
					return err
				}
			}
			return nil
		})
	}
	return nil
}
