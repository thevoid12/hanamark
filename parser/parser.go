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
	"sort"
	"strings"

	"github.com/spf13/viper"
)

// root files need to be updated in the end after all parsing in dest folder is done
func ParseFiles(ctx context.Context) error {

	l := logs.GetLoggerctx(ctx)

	// baseFileMap has the list of base files few has subfolder in them so these basefiles
	// should have links to the subfiles in it eg blogs.html should have links of all the blogs
	// in its corresponding subfolder few are just md files that needs to be parsed
	baseFileMap := viper.GetStringMapString("fileMeta.baseFilesMap")

	var metaList []*model.PageMeta

	for basefileName, bfdir := range baseFileMap {
		ext := filepath.Ext(bfdir)
		if ext != ".md" { // ie it is a directory so we go parse sub directory
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
			metaList, err = parseSubFolderFilesToHtml(ctx, bfdir)
			if err != nil {
				l.Sugar().Error("parse subfolder files to html failed", err)
				return err
			}
			// since all the files in the subfolder is parsed we will now process the index page for these subfolder(base file)
			// of if there are no sub folder the base file md is directly converted to html
			err = tmplt.RenderBaseLinkTemplate(ctx, metaList, basefileName)
			if err != nil {
				return err
			}

		} else {
			// process pure base files(which has no subdirectory)
			rootSrcDir := viper.GetString("filepath.sourceMDRoot")
			fp := filepath.Join(rootSrcDir, bfdir)
			info, err := os.Stat(fp)
			if err != nil {
				l.Sugar().Error("src file not found", err)
				return err
			}
			meta, err := parseMarkDownFile(ctx, fp, basefileName, info)
			if err != nil {
				return err
			}
			metaList = append(metaList, meta)
		}

	}
	return nil
}

func parseSubFolderFilesToHtml(ctx context.Context, baseFiledir string) (metaList []*model.PageMeta, err error) {
	rootSrcDir := viper.GetString("filepath.sourceMDRoot")

	// traverse through the sub directory of src  and create links to the base file in destination
	err = filepath.Walk(filepath.Join(rootSrcDir, baseFiledir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only Markdown files
		meta, err := parseMarkDownFile(ctx, path, baseFiledir, info)
		if err != nil {
			return err
		}
		if meta != nil {
			metaList = append(metaList, meta)
		}
		return nil
	})

	if len(metaList) > 1 {
		// Sorting based on Date field in desc order so that latest record is always at the top
		sort.SliceStable(metaList, func(i, j int) bool {
			return metaList[i].Date.After(metaList[j].Date)
		})
	}
	return metaList, err
}

func parseMarkDownFile(ctx context.Context, path, baseFiledir string, info os.FileInfo) (meta *model.PageMeta, err error) {
	l := logs.GetLoggerctx(ctx)

	rootSrcDir := viper.GetString("filepath.sourceMDRoot")
	rootDestDir := viper.GetString("filepath.destMDRoot")

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
		// Determine relative path from source root
		relPath, err := filepath.Rel(rootSrcDir, path)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		lastModfiedTime := info.ModTime()
		// Generate markdown with file links
		GeneratedHtml, err := ParseMarkdownToHtml(path)
		if err != nil {
			l.Sugar().Error("Error parsing markdown to html", err)
			return nil, err
		}
		title := ""
		if GeneratedHtml != "" {
			title, err = ExtractHeadingInMarkdown(ctx, path)
			if err != nil {
				return nil, err
			}
		}
		meta = &model.PageMeta{
			GenHtml:     GeneratedHtml,
			PageName:    "",
			PageTitle:   title,
			Date:        lastModfiedTime,
			DestPageDir: destPath,
			BaseFile:    baseFiledir,
		}
		outputHtml, err := tmplt.RenderTemplate(ctx, meta)
		if err != nil {
			return nil, err
		}
		err = tmplt.WriteIntoFile(ctx, outputHtml, meta)
		if err != nil {
			return nil, err
		}
		meta.GenHtml = "" // there is no use of storing it in memory
		destPath = util.RemoveRootPartOfDir(destPath, viper.GetString("filepath.destMDRoot"))
		meta.DestPageDir = destPath
	}

	return meta, nil
}
