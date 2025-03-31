package parser

import (
	"context"
	"errors"
	logs "hanamark/logger"
	"hanamark/model"
	tmplt "hanamark/templates"
	"hanamark/util"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// root files need to be updated in the end after all parsing in dest folder is done
func SaveBasefile(ctx context.Context) error {

	l := logs.GetLoggerctx(ctx)

	// baseFileMap has the list of base files which has subfolder in them so these basefiles
	// should have links to the subfiles in it eg blogs.html should have links of all the blogs
	// in its corresponding subfolder
	baseFileMap := viper.GetStringMapString("fileMeta.baseFilesMap")

	for basefileName, bfdir := range baseFileMap {

		rootDestDir := viper.GetString("filepath.destMDRoot")
		newdir := filepath.Join(rootDestDir, bfdir)
		_, err := os.Stat(newdir)
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(newdir, 0755)
			if err != nil {
				l.Sugar().Error("create file failed", err)
				return err
			}
		} else if err != nil {
			l.Sugar().Error("file path not found", err)
			return err
		}
		metaList, err := parseSubFolderFilesToHtml(ctx, bfdir)
		if err != nil {
			l.Sugar().Error("parse subfolder files to html failed", err)
			return err
		}

		// since all the files in the subfolder is parsed we will now process the index page for these subfolder(base file)

		err = tmplt.RenderBaseTemplate(ctx, metaList, basefileName)
		if err != nil {
			return err
		}

	}
	return nil
}

func parseSubFolderFilesToHtml(ctx context.Context, baseFiledir string) (metaList []*model.PageMeta, err error) {
	l := logs.GetLoggerctx(ctx)

	rootSrcDir := viper.GetString("filepath.sourceMDRoot")
	rootDestDir := viper.GetString("filepath.destMDRoot")

	// traverse through the sub directory of src  and create links to the base file in destination
	err = filepath.Walk(filepath.Join(rootSrcDir, baseFiledir), func(path string, info os.FileInfo, err error) error {
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
			destPath = util.RemoveExtentionFromFile(destPath)
			destPath += ".html"
			destDir := filepath.Dir(destPath)
			// Ensure the destination directory exists
			err = os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				l.Sugar().Error("make destination director failed", err)
				return err
			}

			// Generate markdown with file links
			GeneratedHtml, err := ParseMarkdownToHtml(path)
			if err != nil {
				l.Sugar().Error("Error parsing markdown to html", err)
				return err
			}
			tilte, err := ExtractHeadingInMarkdown(ctx, path)
			if err != nil {
				return err
			}
			meta := &model.PageMeta{
				GenHtml:     GeneratedHtml,
				PageName:    "",
				PageTitle:   tilte,
				Date:        time.Now(),
				DestPageDir: destPath,
				PageType:    "",
				BaseFile:    baseFiledir,
			}
			outputHtml, err := tmplt.RenderTemplate(ctx, meta)
			if err != nil {
				return err
			}
			err = tmplt.WriteIntoFile(ctx, outputHtml, meta)
			if err != nil {
				return err
			}
			meta.GenHtml = "" // there is no use of storing it in memory
			destPath = util.RemoveRootPartOfDir(destPath, viper.GetString("filepath.destMDRoot"))
			meta.DestPageDir = destPath
			metaList = append(metaList, meta)
		}
		return nil
	})
	return metaList, err
}
