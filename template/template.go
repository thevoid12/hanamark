package tmplt

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"hanamark/util"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

func getTemplate(targetTemplatePath string) (*template.Template, string, error) {
	// l := logs.GetLoggerctx(ctx)
	templateRoot := viper.GetString("filepath.templatePath")
	basePath := filepath.Join(templateRoot, "_base.html")

	var err error

	// Check if _base.html exists
	useBase := false
	if _, statErr := os.Stat(basePath); statErr == nil {
		useBase = true
	}

	// Always parse the target template content first (with define "main")
	// because we want to define the "main" block for the base template to use.
	content, err := os.ReadFile(targetTemplatePath)
	if err != nil {
		return nil, "", err
	}

	finalContent := string(content)
	if useBase {
		if !strings.Contains(finalContent, "define \"main\"") {
			finalContent = fmt.Sprintf(`{{ define "main" }}%s{{ end }}`, finalContent)
		}
	}

	//  custom template functions
	funcMap := template.FuncMap{
		"config": func(key string) any {
			return viper.Get(key)
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"urlJoin": func(base string, elem ...string) string {
			result, err := url.JoinPath(base, elem...)
			if err != nil {
				return base // fallback or panic, but fallback safer for template
			}
			return result
		},
		"cleanUrl": util.CleanURLPath,
	}

	// Parse the target content first with custom functions
	tmpl := template.New("").Funcs(funcMap)
	tmpl, err = tmpl.Parse(finalContent)
	if err != nil {
		return nil, "", err
	}

	// Then parse all templates from the root directory to support partials
	if entries, readErr := os.ReadDir(templateRoot); readErr == nil {
		var templateFiles []string
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
				path := filepath.Join(templateRoot, entry.Name())
				// Avoid adding _base.html and target template again
				if path != basePath && path != targetTemplatePath {
					templateFiles = append(templateFiles, path)
				}
			}
		}
		if len(templateFiles) > 0 {
			tmpl, err = tmpl.ParseFiles(templateFiles...)
			if err != nil {
				return nil, "", err
			}
		}
	}

	// Then parse base template if it exists
	if useBase {
		// We use ParseFiles on the existing tmpl to add the base template
		tmpl, err = tmpl.ParseFiles(basePath)
		if err != nil {
			return nil, "", err
		}
	}

	// Execution name: if using base, we usually execute "_base.html" (or "base").
	// If base defines 'block "main" .', and target 'define "main"',
	// executing "_base.html" renders the shell with the target's main.
	execName := ""
	if useBase {
		execName = "_base.html"
	}
	return tmpl, execName, nil
}

// takes in the base template and appends the content the base template and gives us back the final html string
func RenderTemplate(ctx context.Context, meta *model.PageMeta, templatePath string, opFile string) error {
	l := logs.GetLoggerctx(ctx)

	tmpl, execName, err := getTemplate(templatePath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	f, err := os.Create(opFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, execName, meta)
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return err
	}

	return nil
}

// RenderBaseLinkTemplate processes files which has links of other sub files. ususlly every folder has individual files
// and the files are linked if it is a sub folder!
func RenderBaseLinkTemplate(ctx context.Context, metaList []*model.PageMeta, lp *model.ListPage) error {
	l := logs.GetLoggerctx(ctx)
	baseFolderName := lp.Base
	tmptPath := lp.TempPath
	base := filepath.Base(baseFolderName)

	opBaseFile := filepath.Join(viper.GetString("filepath.destHtmlDir"), baseFolderName, "index.html")

	if len(metaList) > 1 {
		// Sorting based on Date field in desc order so that latest record is always at the top
		isSortFileByCreatedOn := viper.GetBool("sortFilesByCreatedOn")
		if isSortFileByCreatedOn {
			sort.SliceStable(metaList, func(i, j int) bool {
				return metaList[i].CreatedDate.After(metaList[j].CreatedDate)
			})
		} else { // sort by updated date
			sort.SliceStable(metaList, func(i, j int) bool {
				return metaList[i].UpdatedDate.After(metaList[j].UpdatedDate)
			})
		}
	}

	tmpl, execName, err := getTemplate(tmptPath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	f, err := os.Create(opBaseFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()

	// Get the directory where the list page (index.html) will be created
	listPageDir := filepath.Dir(opBaseFile)

	for _, m := range metaList {
		// DestPageDir at this point is like "./blogs/2026_01/updates_01.html"
		// We need to calculate the relative path from the list page directory to the target file

		// First, strip the leading "./" to get an absolute path relative to destHtmlDir
		destPath := strings.TrimPrefix(m.DestPageDir, "./")

		// Build the full path
		fullDestPath := filepath.Join(viper.GetString("filepath.destHtmlDir"), destPath)

		// Calculate relative path from list page directory to the target file
		dir, relErr := filepath.Rel(listPageDir, fullDestPath)
		if relErr != nil {
			l.Sugar().Error("error in getting relative path", relErr)
			return relErr
		}
		m.DestPageDir = dir
	}
	// wrap data for base template
	pageTitle := base

	// Check if this is the root index.html and if titletag is configured
	if baseFolderName == "." {
		titleTag := viper.GetString("indexHomepageHtml.titletag")
		if titleTag != "" {
			pageTitle = titleTag
		}
	}

	data := map[string]interface{}{
		"PageTitle": pageTitle,
		"List":      metaList,
	}

	err = tmpl.ExecuteTemplate(f, execName, data)
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return err
	}

	return nil
}

func RenderBaseTagListTemplate(ctx context.Context, taglist []*model.TagList, tmptPath string) error {
	l := logs.GetLoggerctx(ctx)

	// TODO: if there is a filemeta to change the base folder name give it more preced

	opBaseFile := filepath.Join(viper.GetString("filepath.destHtmlDir"), "tags", "index.html")

	tmpl, execName, err := getTemplate(tmptPath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	dir := filepath.Dir(opBaseFile)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		l.Sugar().Error("failed to create directories", err)
		return err
	}
	f, err := os.Create(opBaseFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()

	// wrap data for base template
	data := map[string]interface{}{
		"PageTitle": "Tags",
		"List":      taglist,
	}

	err = tmpl.ExecuteTemplate(f, execName, data)
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return err
	}

	return nil
}

func RenderTagLinkTemplate(ctx context.Context, tagMeta []*model.Tag, tagName string) error {
	l := logs.GetLoggerctx(ctx)
	baseFolderName := tagMeta[0].TagDestPath

	// TODO: if there is a filemeta to change the base folder name give it more preced

	// opBaseFile := filepath.Join(viper.GetString("filepath.destHtmlDir"), baseFolderName, strings.TrimSuffix(baseFolderName, filepath.Ext(baseFolderName))+".html")
	opBaseFile := baseFolderName

	tmpl, execName, err := getTemplate(tagMeta[0].TagTemplatePath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	dir := filepath.Dir(opBaseFile)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		l.Sugar().Error("failed to create directories", err)
		return err
	}
	f, err := os.Create(opBaseFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()

	// wrap data for base template
	data := map[string]interface{}{
		"PageTitle": tagName,
		"List":      tagMeta,
	}

	err = tmpl.ExecuteTemplate(f, execName, data)
	if err != nil {
		l.Sugar().Error("Error executing tag meta template", err)
		return err
	}

	return nil
}
