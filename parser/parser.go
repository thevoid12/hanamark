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
	"time"

	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

// root files need to be updated in the end after all parsing in dest folder is done
func ParseFiles(ctx context.Context) error {

	l := logs.GetLoggerctx(ctx)

	sourceFilePath := viper.GetString("filepath.sourceMDDir")
	if sourceFilePath == "" {
		return errors.New("sourceMDDir is empty")
	}
	destRootPath := viper.GetString("filepath.destHtmlDir")
	if destRootPath == "" {
		return errors.New("dest root path in config is empty")
	}
	templateRootPath := viper.GetString("filepath.templatePath")
	if templateRootPath == "" {
		return errors.New("templatePath is empty")
	}

	folderMetaMap := make(map[string][]*model.PageMeta) // key is the folder
	indexFmMap := make(map[string]model.FrontMatter)    // front matter of _index.md page
	var listPages []*model.ListPage                     // this has the base file names (folder names) of the pages
	tagMap := make(map[string][]*model.Tag)             // key is the tag name value is the tag property
	indexMdExists := false                              // track if index.md exists in root

	// trying to implement a mirror tree walker
	err := filepath.WalkDir(sourceFilePath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr // handle permission errors etc.
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
		assetsRel, err := filepath.Rel(sourceFilePath, viper.GetString("filepath.mdAssetsSourcePath"))
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

		// Check if index.md exists in root
		if !d.IsDir() && d.Name() == "index.md" && relSource == "index.md" {
			indexMdExists = true
		}

		if !d.IsDir() && d.Name() == "_index.md" {
			dirRel := filepath.Dir(relSource)
			templatePath = filepath.Join(templateRootPath, dirRel)
			fm, err := ParseFrontMatter(ctx, path)
			if err != nil {
				return err
			}
			relSourcePath, err := filepath.Rel(viper.GetString("filepath.sourceMDDir"), path)
			if err != nil {
				return err
			}
			folderName := filepath.Dir(relSourcePath)
			if _, ok := indexFmMap[folderName]; !ok && len(fm) > 0 {
				indexFmMap[folderName] = fm
			}
			templatefm := "list.html"
			isPresent, val, err := FrontMatterValidator(ctx, path, fm, model.TEMPLATE)
			if err != nil {
				return err
			}
			if isPresent {
				if s, ok := val.(string); ok {
					templatefm = s
				} else {
					l.Sugar().Errorf("template value in frontmatter is not a string for path: %s", path)
				}
			}

			listTemplate, err = util.FindTemplateUpward(templatePath, templateRootPath, templatefm)
			if err != nil {
				return fmt.Errorf("failed to find %s template for directory %s: %w", templatefm, templatePath, err)
			}
			fmt.Println("B:", listTemplate)
			l.Info("B:" + listTemplate)

			listPages = append(listPages, &model.ListPage{Base: folderName, TempPath: listTemplate})
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
				var val interface{}
				isPresent, val, err = FrontMatterValidator(ctx, path, fm, model.DRAFT)
				if err != nil {
					return err
				}
				if isPresent {
					if isDraft, ok := val.(bool); ok && isDraft {
						return nil // skip the file
					} else if !ok {
						l.Sugar().Errorf("draft value in frontmatter is not a boolean for path: %s", path)
					}
				}
			}
			templatefm := ""
			var tags []string
			if len(fm) != 0 {
				var val interface{}
				isPresent, val, err = FrontMatterValidator(ctx, path, fm, model.TEMPLATE)
				if err != nil {
					return err
				}
				if isPresent {
					if s, ok := val.(string); ok {
						templatefm = s
					} else {
						l.Sugar().Errorf("template value in frontmatter is not a string for path: %s", path)
					}
				}

				isPresent, val, err = FrontMatterValidator(ctx, path, fm, model.TAGS)
				if err != nil {
					return err
				}
				if isPresent {
					raw, ok := val.([]interface{})
					if !ok {
						l.Sugar().Errorf("tags value in frontmatter is not an array for path: %s", path)
					} else {
						for i, v := range raw {
							s, typeOk := v.(string)
							if !typeOk {
								return fmt.Errorf("tag at index %d is not a string", i)
							}
							tags = append(tags, s)
						}
					}
				}
			}
			if templatefm == "" {
				// Search for single.html starting from templatePath and going upward to templateRootPath
				singleTemplate, err = util.FindTemplateUpward(templatePath, templateRootPath, "single.html")
				if err != nil {
					return fmt.Errorf("failed to find single.html template for directory %s: %w", templatePath, err)
				}
			} else {
				singleTemplate = filepath.Join(templateRootPath, templatefm)
				// Verify custom template exists
				if info, statErr := os.Stat(singleTemplate); errors.Is(statErr, os.ErrNotExist) || info.IsDir() {
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
					tagMeta, tagErr := ProcessTags(ctx, tag)
					if tagErr != nil {
						return tagErr
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

			relSourcePath, err := filepath.Rel(viper.GetString("filepath.sourceMDDir"), path)
			if err != nil {
				return err
			}
			folderName := filepath.Dir(relSourcePath)
			if folderMetaMap[folderName] == nil {
				folderMetaMap[folderName] = make([]*model.PageMeta, 0)
			}
			folderMetaMap[folderName] = append(folderMetaMap[folderName], meta)
		}

		return nil
	})

	if err != nil {
		return err
	}
	// parse the list template
	rssFeedItems := []*feeds.Item{}
	newfolderMetaMap := make(map[string][]*model.PageMeta) // key is the folder

	// distribute files and sub-list-pages to their closest parent list page
	for _, lp := range listPages {
		newfolderMetaMap[lp.Base] = make([]*model.PageMeta, 0)
	}

	//  Distribute files to their closest parent list page
	for folderPath, files := range folderMetaMap {
		closestLP := getClosestListPage(folderPath, listPages)
		if closestLP != nil {
			newfolderMetaMap[closestLP.Base] = append(newfolderMetaMap[closestLP.Base], files...)
		}
	}

	//  Add sub-list-pages to their parent list pages as entries
	for _, lp := range listPages {
		if lp.Base == "." || lp.Base == "" {
			continue
		}
		parentDir := filepath.Dir(lp.Base)
		parentLP := getClosestListPage(parentDir, listPages)
		if parentLP != nil && parentLP.Base != lp.Base {
			childBase := filepath.Base(lp.Base)
			destPath := filepath.Join(lp.Base, "index.html")

			lpMeta := &model.PageMeta{
				PageTitle:   childBase,
				DestPageDir: destPath,
			}
			newfolderMetaMap[parentLP.Base] = append(newfolderMetaMap[parentLP.Base], lpMeta)
		}
	}

	indexHomepageType := strings.ToLower(strings.TrimSpace(viper.GetString("indexHomepageHtml.type")))
	indexHomepageName := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(viper.GetString("indexHomepageHtml.name")), "/"))

	for _, lp := range listPages {
		if newfolderMetaMap[lp.Base] == nil {
			return errors.New("no files found in the directory:" + lp.Base)
		}
		_, ok := indexFmMap[lp.Base]
		if ok {
			// get the custom template name from the map
			var isPresent bool
			var val interface{}
			isPresent, val, err = FrontMatterValidator(ctx, "", indexFmMap[lp.Base], model.TEMPLATE)
			if err != nil {
				return err
			}
			if isPresent {
				if s, castOk := val.(string); castOk {
					lp.TempPath = filepath.Join(templateRootPath, s) // custom template instead of list.html
				} else {
					l.Sugar().Errorf("template value in _index.md frontmatter is not a string for list page: %s", lp.Base)
				}
			}

			isRssEnabled := viper.GetBool("rss.isRssEnabled")
			if isRssEnabled {
				isPresent, _, err = FrontMatterValidator(ctx, "", indexFmMap[lp.Base], model.RSS)
				if err != nil {
					return err
				}
				if isPresent { // get rss feed items
					items, rssErr := GetRssFeedItems(newfolderMetaMap[lp.Base])
					if rssErr != nil {
						return rssErr
					}
					rssFeedItems = append(rssFeedItems, items...)
				}
			}
		}
		err = tmplt.RenderBaseLinkTemplate(ctx, newfolderMetaMap[lp.Base], lp)
		if err != nil {
			return err
		}

		// If no index.md exists and indexHomepageHtml type is section and matches this section, render to root index.html
		if !indexMdExists && indexHomepageType == model.IndexTypeSection && indexHomepageName != "" && strings.ToLower(strings.TrimSuffix(lp.Base, "/")) == indexHomepageName {
			// Create a copy of lp for root index.html
			rootLp := &model.ListPage{
				Base:     ".",
				TempPath: lp.TempPath,
			}

			// check if _index.html exists
			_, err = os.Stat("_index.html")
			if err != nil && os.IsNotExist(err) {
				// there is no custom index template so
				// we will use the reference of the referenced template section
				rootLp.TempPath = lp.TempPath
			} else if err != nil {
				return err
			} else {
				// Use _index.html template if exists, otherwise use the section's template
				var indexTemplate string
				indexTemplate, err = util.FindTemplateUpward(templateRootPath, templateRootPath, "_index.html")
				if err == nil {
					rootLp.TempPath = indexTemplate
				}
			}

			rootPageMeta := make([]*model.PageMeta, len(newfolderMetaMap[lp.Base]))
			for i, pm := range newfolderMetaMap[lp.Base] {
				adjustedPm := *pm
				adjustedPm.DestPageDir = filepath.Join(lp.Base, pm.DestPageDir)
				rootPageMeta[i] = &adjustedPm
			}
			err = tmplt.RenderBaseLinkTemplate(ctx, rootPageMeta, rootLp)
			if err != nil {
				return err
			}
			l.Info("index.html generated from indexHomepageHtml section: " + indexHomepageName)
		}
	}

	// If no index.md and indexHomepageHtml type is page, copy that page to index.html
	if !indexMdExists && indexHomepageType == model.IndexTypePage && indexHomepageName != "" {
		// Normalize to .html if needed
		pageName := indexHomepageName
		if !strings.HasSuffix(pageName, ".html") {
			pageName += ".html"
		}
		// Copy the rendered page to index.html
		srcPath := filepath.Join(destRootPath, pageName)
		destPath := filepath.Join(destRootPath, "index.html")
		if copyErr := util.CopyFile(srcPath, destPath); copyErr != nil {
			l.Sugar().Warn("indexHomepageHtml page not found, skipping: " + pageName)
		} else {
			l.Info("index.html generated from indexHomepageHtml page: " + pageName)
		}
	}
	// if rss is enabled in Frontmatter we got to render that
	if len(rssFeedItems) > 0 {
		// parse rss for the list
		var feed *feeds.Feed
		feed, err = GetRssFeed()
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
		err = tmplt.RenderTagLinkTemplate(ctx, tagMeta, tagName)
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
	isTagsEnabled := viper.GetBool("tags")
	if isTagsEnabled {
		tagListTmpt := filepath.Join(templateRootPath, "tags", "single.html")
		if tagErr := tmplt.RenderBaseTagListTemplate(ctx, tagList, tagListTmpt); tagErr != nil {
			l.Sugar().Errorf("failed to render tag list template: %v", tagErr)
		}
	}
	return nil
}

func processFile(ctx context.Context, sourcePath string, templatePath string, fm map[string]any) (*model.PageMeta, error) {

	l := logs.GetLoggerctx(ctx)
	// check if the parent folders exists. if not create the parent folders
	basefileName := filepath.Base(sourcePath) //TODO: this doesnt work as base includes just the last file name but we want everything other than the root to mirror destination
	rootDestDir := viper.GetString("filepath.destHtmlDir")
	if rootDestDir == "" {
		return nil, errors.New("destination root directory is not set")
	}
	relSourcePath, err := filepath.Rel(viper.GetString("filepath.sourceMDDir"), sourcePath)
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

	info, err := os.Stat(sourcePath)
	if err != nil {
		l.Sugar().Error("src file not found", err)
		return nil, err
	}
	return parseMarkDownFile(ctx, sourcePath, basefileName, info, templatePath, fm)
}

func parseMarkDownFile(ctx context.Context, path, baseFiledir string, info os.FileInfo, templatePath string, fm map[string]any) (meta *model.PageMeta, err error) {
	l := logs.GetLoggerctx(ctx)

	rootSrcDir := viper.GetString("filepath.sourceMDDir")
	rootDestDir := viper.GetString("filepath.destHtmlDir")

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
		// Determine relative path from source root
		var relPath string
		relPath, err = filepath.Rel(rootSrcDir, path) // TODO: repeating twice
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

		isPresent, val, err := FrontMatterValidator(ctx, path, fm, model.DATE)
		if err != nil {
			return nil, err
		}
		if !isPresent {
			return nil, errors.New("cannot render file without created_on data in frontmatter")
		}
		var createdOn time.Time
		if s, ok := val.(string); ok {
			createdOn, err = util.ParseTimeFlexible(s)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("created_on value in frontmatter is not a string for path: %s", path)
		}

		lastModfiedTime := info.ModTime()
		isPresentUpdated, updatedVal, err := FrontMatterValidator(ctx, path, fm, model.UPDATED_ON)
		if err != nil {
			return nil, err
		}
		if isPresentUpdated {
			if s, ok := updatedVal.(string); ok {
				t, err := util.ParseTimeFlexible(s)
				if err != nil {
					return nil, err
				}
				lastModfiedTime = t
			}
		}

		// Generate markdown with file links
		generatedHtml, err := ParseMarkdownToHtml(ctx, path)
		if err != nil {
			l.Sugar().Error("Error parsing markdown to html", err)
			return nil, err
		}
		isPresent, value, err := FrontMatterValidator(ctx, "", fm, model.TITLE)
		if err != nil {
			return nil, err
		}
		title := ""
		if isPresent {
			if s, ok := value.(string); ok {
				title = s
			}
		}
		if generatedHtml != "" {
			if len(fm) > 0 {
				// we have to remove the frontmatter from the md else it will also be displayed along the html
				generatedHtml = string(StripFrontMatter(ctx, []byte(generatedHtml)))
			}
			if !isPresent {
				title, err = ExtractHeadingInMarkdown(ctx, path)
				if err != nil {
					return nil, err
				}
			}
		}

		// Calculate read time from markdown content
		mdContent, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		// Strip frontmatter before calculating read time
		if len(fm) > 0 {
			mdContent = StripFrontMatter(ctx, mdContent)
		}
		readTime := util.CalculateReadTime(string(mdContent))

		diskPath := destPath
		destPath = util.RemoveRootPartOfDir(destPath, viper.GetString("filepath.destHtmlDir"))

		meta = &model.PageMeta{
			GenHtml:        generatedHtml,
			ReadTime:       readTime,
			PageName:       "",
			PageTitle:      title,
			CreatedDate:    createdOn,
			UpdatedDate:    lastModfiedTime,
			DestPageDir:    destPath,
			BaseFile:       baseFiledir,
			FrontMatterMap: fm,
		}
		err = tmplt.RenderTemplate(ctx, meta, templatePath, diskPath)
		if err != nil {
			return nil, err
		}
	}

	return meta, nil
}

func getClosestListPage(targetPath string, listPages []*model.ListPage) *model.ListPage {
	var closest *model.ListPage
	for _, lp := range listPages {
		if util.IsDirUnder(lp.Base, targetPath) {
			if closest == nil || len(lp.Base) > len(closest.Base) {
				closest = lp
			}
		}
	}
	return closest
}
