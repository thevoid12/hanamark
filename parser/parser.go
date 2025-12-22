package parser

import (
	"context"
	"errors"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	tmplt "hanamark/template"
	"hanamark/util"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

// root files need to be updated in the end after all parsing in dest folder is done
func ParseFiles(ctx context.Context) error {

	l := logs.GetLoggerctx(ctx)

	sourceFilePath := viper.GetString("filepath.sourceMDRoot")
	if sourceFilePath == "" {
		return errors.New("sourceMDRoot is empty")
	}
	destRootPath := viper.GetString("filepath.destHtmlRoot")
	if destRootPath == "" {
		return errors.New("dest root path in config is empty")
	}
	templateRootPath := viper.GetString("Filepath.templatePath")
	if templateRootPath == "" {
		return errors.New("templatePath is empty")
	}

	folderMetaMap := make(map[string][]*model.PageMeta) // key is the folder
	indexFmMap := make(map[string]model.FrontMatter)    // front matter of _index.md page
	var ListPages []*model.ListPage                     // this has the base file names (folder names) of the pages
	tagMap := make(map[string][]*model.Tag)             // key is the tag name value is the tag property

	// trying to implement a mirror tree walker
	err := filepath.WalkDir(sourceFilePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // handle permission errors etc.
		}

		base := filepath.Base(path)
		l.Info(path)
		relSource, err := filepath.Rel(sourceFilePath, path)
		if err != nil {
			return err
		}

		// Skip hidden files and directories (starting with '.')
		if strings.HasPrefix(base, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil // skip hidden files
		}

		// continue or ignore if we see assest folder. assets folder has noting to do with parsing
		// fetching the relative path for assets
		assetsRel, err := filepath.Rel(sourceFilePath, viper.GetString("filepath.sourceAssetsPath"))
		if err != nil {
			return err
		}
		// Skip assets directory by checking both the config path and the directory name
		if d.IsDir() && (relSource == assetsRel) {
			return filepath.SkipDir
		}

		// Corresponding path in template
		templatePath := filepath.Join(templateRootPath, relSource)
		// if the path is a file then we need to get the parent directory
		singleTemplate := ""
		listTemplate := ""
		if !d.IsDir() && d.Name() == "_index.md" {
			dirRel := filepath.Dir(relSource)
			templatePath = filepath.Join(templateRootPath, dirRel)
			fm, err := ParseFrontMatter(ctx, path)
			if err != nil {
				return err
			}
			relSourcePath, err := filepath.Rel(viper.GetString("filepath.sourceMDRoot"), path)
			if err != nil {
				return err
			}
			folderName := filepath.Dir(relSourcePath)
			if indexFmMap[folderName] == nil && fm != nil {
				indexFmMap[folderName] = fm
			}
			templatefm := "list.html"
			isPresent, val, err := FrontMatterValidator(ctx, path, fm, model.TEMPLATE)
			if err != nil {
				return err
			}
			if isPresent {
				templatefm = val.(string)
			}

			listTemplate, err = util.FindTemplateUpward(templatePath, templateRootPath, templatefm)
			if err != nil {
				return fmt.Errorf("failed to find %s template for directory %s: %w", templatefm, templatePath, err)
			}
			fmt.Println("B:", listTemplate)
			l.Info("B:" + listTemplate)

			ListPages = append(ListPages, &model.ListPage{Base: folderName, TempPath: listTemplate})
			return nil
		} else if !d.IsDir() {
			// Skip non-markdown files
			if !strings.HasSuffix(base, ".md") {
				return nil
			}
			templatePath = filepath.Dir(templatePath)
			// non directory will have single templates!
			fm, err := ParseFrontMatter(ctx, path)
			if err != nil {
				return err
			}

			isPresent, _, err := FrontMatterValidator(ctx, path, fm, model.DATE)
			if err != nil {
				return err
			}
			if !isPresent {
				return errors.New("cannot render file without created_on data in frontmatter.file:" + path)
			}
			// Check if the file is in draft state. If so we can ignore the file and move forward
			if len(fm) != 0 {
				isPresent, val, err := FrontMatterValidator(ctx, path, fm, model.DRAFT)
				if err != nil {
					return err
				}
				if isPresent && val.(bool) {
					return nil // skip the file
				}
			}
			templatefm := ""
			var tags []string
			if len(fm) != 0 {
				isPresent, val, err := FrontMatterValidator(ctx, path, fm, model.TEMPLATE)
				if err != nil {
					return err
				}
				if isPresent {
					templatefm = val.(string)
				}

				isPresent, val, err = FrontMatterValidator(ctx, path, fm, model.TAGS)
				if err != nil {
					return err
				}
				if isPresent {
					raw := val.([]interface{})
					for i, v := range raw {
						s, ok := v.(string)
						if !ok {
							return fmt.Errorf("tag at index %d is not a string", i)
						}
						tags = append(tags, s)
					}
				}

			}
			if templatefm == "" {
				// Search for single.html starting from templatePath and going upward to templateRootPath
				var err error
				singleTemplate, err = util.FindTemplateUpward(templatePath, templateRootPath, "single.html")
				if err != nil {
					return fmt.Errorf("failed to find single.html template for directory %s: %w", templatePath, err)
				}
			} else {
				singleTemplate = filepath.Join(templateRootPath, templatefm)
				// Verify custom template exists
				if info, err := os.Stat(singleTemplate); errors.Is(err, os.ErrNotExist) || info.IsDir() {
					return fmt.Errorf("custom template %s is missing or is a directory", singleTemplate)
				} else if err != nil {
					return err
				}
			}
			fmt.Println("A:", singleTemplate)
			l.Info("A:" + singleTemplate)

			meta, err := processFile(ctx, path, singleTemplate, fm)
			if err != nil {
				return err
			}
			// TODO: parse tag
			// process tag
			if len(tags) > 0 {
				for _, tag := range tags {
					tagMeta, err := ProcessTags(ctx, tag)
					if err != nil {
						return err
					}

					if tagMap[tag] == nil {
						tagMap[tag] = make([]*model.Tag, 0)
					}
					tagMeta.FileHeading = meta.PageTitle

					destPath := filepath.Join(destRootPath, meta.DestPageDir)
					relDir, err := util.RelURL(tagMeta.TagDestPath, destPath)
					if err != nil {
						return err
					}
					tagMeta.FileDestPath = relDir
					tagMap[tag] = append(tagMap[tag], tagMeta)
				}
			}

			meta.FrontMatterMap = fm

			relSourcePath, err := filepath.Rel(viper.GetString("filepath.sourceMDRoot"), path)
			if err != nil {
				return err
			}
			folderName := filepath.Dir(relSourcePath)
			if folderMetaMap[folderName] == nil {
				folderMetaMap[folderName] = make([]*model.PageMeta, 0)
			}
			folderMetaMap[folderName] = append(folderMetaMap[folderName], meta)
		}
		// } else {
		// 	// only if all the files in the folders are traversed and we have the metaList we can process the list template coz list is a collection of links to the files in the folder

		// 	// in root directory we dont need list.html as if there is any list definitely there will be a subfolder
		// 	if relSource != "." {
		// 		// Search for list.html starting from templatePath and going upward to templateRootPath
		// 		var err error
		// 		listTemplate, err = util.FindTemplateUpward(templatePath, templateRootPath, "list.html")
		// 		if err != nil {
		// 			return fmt.Errorf("failed to find list.html template for directory %s: %w", templatePath, err)
		// 		}
		// 		fmt.Println("B:", listTemplate)
		// 		l.Info("B:" + listTemplate)
		// 		ListPages = append(ListPages, &model.ListPage{Base: base, TempPath: listTemplate})
		// 		// if metaMap[base] != nil { // TODO: but what if we reach here before processing the files?

		// 		// } else {
		// 		// 	return errors.New("no files found in the directory" + templatePath)
		// 		// }
		// 	}
		// }
		return nil
	})

	if err != nil {
		return err
	}
	// parse the list template
	rssFeedItems := []*feeds.Item{}
	for _, lp := range ListPages {
		if folderMetaMap[lp.Base] == nil {
			return errors.New("no files found in the directory:" + lp.Base)
		}
		if indexFmMap[lp.Base] != nil {
			// get the custom template name from the map
			isPresent, val, err := FrontMatterValidator(ctx, "", indexFmMap[lp.Base], model.TEMPLATE)
			if err != nil {
				return err
			}
			if isPresent {
				lp.TempPath = filepath.Join(templateRootPath, val.(string)) // custom template instead of list.html
			}

			isRssEnabled := viper.GetBool("rss.isRssEnabled")
			if isRssEnabled {
				isPresent, _, err = FrontMatterValidator(ctx, "", indexFmMap[lp.Base], model.RSS)
				if err != nil {
					return err
				}
				if isPresent { // get rss feed items
					items, err := GetRssFeedItems(folderMetaMap[lp.Base])
					if err != nil {
						return err
					}
					rssFeedItems = append(rssFeedItems, items...)
				}
			}
		}
		err := tmplt.RenderBaseLinkTemplate(ctx, folderMetaMap[lp.Base], lp)
		if err != nil {
			return err
		}

	}
	// if rss is enabled in Frontmatter we got to render that
	if len(rssFeedItems) > 0 {
		// parse rss for the list
		feed, err := GetRssFeed()
		if err != nil {
			return err
		}
		feed.Items = rssFeedItems
		err = GenerateRss(ctx, feed)
		if err != nil {
			return err
		}
	}
	// parse the tag list pages (1 html page for each tag which has the list of stuff thats been tagged)
	var tagList []*model.TagList // tag list is the data needed to create the tags page which will have the list of tags
	for tagName, tagMeta := range tagMap {
		err := tmplt.RenderTagLinkTemplate(ctx, tagMeta, tagName)
		if err != nil {
			return err
		}
		tagDest, err := util.RelURL(filepath.Join(destRootPath, "tags"), tagMeta[0].TagDestPath)
		if err != nil {
			return err
		}
		tagList = append(tagList, &model.TagList{
			TagName:     tagName,
			TagDestPath: tagDest,
			Count:       len(tagMeta),
		})
	}
	// render tag list page (TODO: you should do all of these only if config is set true)
	tagListTmpt := filepath.Join(templateRootPath, "tags", "single.html")
	tmplt.RenderBaseTagListTemplate(ctx, tagList, tagListTmpt)
	return nil
}

func processFile(ctx context.Context, sourcePath string, templatePath string, fm map[string]any) (*model.PageMeta, error) {
	// sp: ./pointA/about.md
	//tp: ./templates/single.html
	// result: ./pointB/about.html
	// config: ./pointB
	//  ./pointB/about.html
	//  sp: ./pointA/home/blog1.md
	//  tp: ./templaates/home/single.html
	//  result: ./pointB/home/blog1.html
	l := logs.GetLoggerctx(ctx)
	// ext := filepath.Ext(sourcePath)
	// check if the parent folders exists. if not create the parent folders
	basefileName := filepath.Base(sourcePath) //TODO: this doesnt work as base includes just the last file name but we want everything other than the root to mirror destination
	rootDestDir := viper.GetString("filepath.destHtmlRoot")
	if rootDestDir == "" {
		return nil, errors.New("destination root directory is not set")
	}
	relSourcePath, err := filepath.Rel(viper.GetString("filepath.sourceMDRoot"), sourcePath)
	if err != nil {
		return nil, err
	}
	// relSource:=filepath.Dir(rel)
	newdir := filepath.Join(rootDestDir, filepath.Dir(relSourcePath))
	_, err = os.Stat(newdir)
	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(newdir, 0755)
		if err != nil {
			l.Sugar().Error("create file failed", err)
			return nil, err
		}

	}
	// metaList, err = parseSubFolderFilesToHtml(ctx, bfdir)
	// if err != nil {
	// 	l.Sugar().Error("parse subfolder files to html failed", err)
	// 	return err
	// }
	// since all the files in the subfolder is parsed we will now process the index page for these subfolder(base file)
	// of if there are no sub folder the base file md is directly converted to html
	// err = tmplt.RenderBaseLinkTemplate(ctx, metaList, basefileName)
	// if err != nil {
	// 	return err
	// }

	// } else {
	// process pure base files(which has no subdirectory)
	// rootSrcDir := viper.GetString("filepath.sourceMDRoot")
	// fp := filepath.Join(rootSrcDir, bfdir)
	info, err := os.Stat(sourcePath)
	if err != nil {
		l.Sugar().Error("src file not found", err)
		return nil, err
	}
	return parseMarkDownFile(ctx, sourcePath, basefileName, info, templatePath, fm)
}

// TODO: this needs to be removed and merged into our new parser mirror tree walker as we are already walking there
// func parseSubFolderFilesToHtml(ctx context.Context, baseFiledir string) (metaList []*model.PageMeta, err error) {
// 	rootSrcDir := viper.GetString("filepath.sourceMDRoot")

// 	// traverse through the sub directory of src  and create links to the base file in destination
// 	err = filepath.Walk(filepath.Join(rootSrcDir, baseFiledir), func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		// Process only Markdown files
// 		meta, err := parseMarkDownFile(ctx, path, baseFiledir, info)
// 		if err != nil {
// 			return err
// 		}
// 		if meta != nil {
// 			metaList = append(metaList, meta)
// 		}
// 		return nil
// 	})

// 	if len(metaList) > 1 {
// 		// Sorting based on Date field in desc order so that latest record is always at the top
// 		sort.SliceStable(metaList, func(i, j int) bool {
// 			return metaList[i].Date.After(metaList[j].Date)
// 		})
// 	}
// 	return metaList, err
// }

func parseMarkDownFile(ctx context.Context, path, baseFiledir string, info os.FileInfo, templatePath string, fm map[string]any) (meta *model.PageMeta, err error) {
	l := logs.GetLoggerctx(ctx)

	rootSrcDir := viper.GetString("filepath.sourceMDRoot")
	rootDestDir := viper.GetString("filepath.destHtmlRoot")

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
		// Determine relative path from source root
		relPath, err := filepath.Rel(rootSrcDir, path) // TODO: repeating twice
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

		isPresent, val, err := FrontMatterValidator(ctx, path, fm, model.DATE)
		if err != nil {
			return nil, err
		}
		if !isPresent {
			return nil, errors.New("cannot render file without created_on data in frontmatter")
		}
		createdOn, err := util.ParseTimeFlexible(val.(string))
		if err != nil {
			return nil, err
		}
		// Generate markdown with file links
		generatedHtml, err := ParseMarkdownToHtml(ctx, path)
		if err != nil {
			l.Sugar().Error("Error parsing markdown to html", err)
			return nil, err
		}
		title := ""
		if generatedHtml != "" {
			if len(fm) > 0 {
				// we have to remove the frontmatter from the md else it will also be displayed along the html
				generatedHtml = string(StripFrontMatter(ctx, []byte(generatedHtml)))
			}
			title, err = ExtractHeadingInMarkdown(ctx, path)
			if err != nil {
				return nil, err
			}
		}

		meta = &model.PageMeta{
			GenHtml:     generatedHtml,
			PageName:    "",
			PageTitle:   title,
			CreatedDate: createdOn,
			UpdatedDate: lastModfiedTime,
			DestPageDir: destPath,
			BaseFile:    baseFiledir,
		}
		err = tmplt.RenderTemplate(ctx, meta, templatePath)
		if err != nil {
			return nil, err
		}
		// err = tmplt.WriteIntoFile(ctx, outputHtml, meta)
		// if err != nil {
		// 	return nil, err
		// }
		destPath = util.RemoveRootPartOfDir(destPath, viper.GetString("filepath.destHtmlRoot"))
		meta.DestPageDir = destPath // TODO: this is bad and this will cause confusion
	}

	return meta, nil
}
