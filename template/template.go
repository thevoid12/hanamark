package tmplt

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

func getTemplate(ctx context.Context, targetTemplatePath string) (*template.Template, string, error) {
	// l := logs.GetLoggerctx(ctx)
	templateRoot := viper.GetString("filepath.templatePath")
	basePath := filepath.Join(templateRoot, "base_.html")

	funcMap := template.FuncMap{
		"config": func(key string) any {
			return viper.Get(key)
		},
		"findAsset": func(assetPath string) string {
			// Extract the filename from the asset path
			filename := filepath.Base(assetPath)
			
			// Get the destination assets directory
			destAssetsPath := viper.GetString("filepath.destAssetsPath")
			if destAssetsPath == "" {
				// Fallback: try to use the original path
				return assetPath
			}
			
			// Walk the destination assets directory to find the file
			var foundPath string
			err := filepath.WalkDir(destAssetsPath, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil // Continue walking even if there's an error
				}
				if !d.IsDir() && filepath.Base(path) == filename {
					// Found the file, make it relative to destMDRoot
					destMDRoot := viper.GetString("filepath.destMDRoot")
					relPath, err := filepath.Rel(destMDRoot, path)
					if err == nil {
						foundPath = "./" + filepath.ToSlash(relPath)
					}
					return filepath.SkipAll // Stop walking once found
				}
				return nil
			})
			
			if err == nil && foundPath != "" {
				return foundPath
			}
			
			// Fallback: return original path
			return assetPath
		},
	}

	tmpl := template.New("").Funcs(funcMap)
	var err error

	// Check if base_.html exists
	useBase := false
	if _, err := os.Stat(basePath); err == nil {
		useBase = true
	}

	filesToParse := []string{}
	// strict ordering: base first, then specific
	if useBase {
		filesToParse = append(filesToParse, basePath)
	}

	// Always parse the target template
	// filesToParse = append(filesToParse, targetTemplatePath)

	// manually reading the file to append {{ define "main" }}
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

	// Also parse partials? For now, we only parse base and target.
	// If the user has separate header.html and base uses it, we should parse it too.
	// But simply ParseFiles(base, target) is the "minimal" approach for "baseof" pattern.
	// We will assume header/footer logic is IN base or inlined for this task,
	// unless we find header.html in the dir, in which case we might want to include "shared" templates.
	// For safety against conflicts, let's just parse base + target.

	if len(filesToParse) > 0 {
		tmpl, err = tmpl.ParseFiles(filesToParse...)
		if err != nil {
			return nil, "", err
		}
	}

	tmpl, err = tmpl.Parse(finalContent)
	if err != nil {
		return nil, "", err
	}

	// Execution name: if using base, we usually execute "base_.html" (or "base").
	// If base defines 'block "main" .', and target 'define "main"',
	// executing "base_.html" renders the shell with the target's main.
	execName := ""
	if useBase {
		execName = "base_.html"
	}
	return tmpl, execName, nil
}

// takes in the base template and appends the content the base template and gives us back the final html string
func RenderTemplate(ctx context.Context, meta *model.PageMeta, templatePath string) error {
	l := logs.GetLoggerctx(ctx)

	tmpl, execName, err := getTemplate(ctx, templatePath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}

	opFile := meta.DestPageDir
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

	tmpl, execName, err := getTemplate(ctx, tmptPath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}

	if len(metaList) > 1 {
		// Sorting based on Date field in desc order so that latest record is always at the top
		isSortFileByCreatedOn := viper.GetBool("sortFileByCreatedOn")
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

	// TODO: if there is a filemeta to change the base folder name give it more preced
	opBaseFile := filepath.Join(viper.GetString("filepath.destMDRoot"), baseFolderName, strings.TrimSuffix(baseFolderName, filepath.Ext(baseFolderName))+".html")
	f, err := os.Create(opBaseFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()
	for _, m := range metaList {
		dir, err := filepath.Rel(baseFolderName, m.DestPageDir)
		if err != nil {
			l.Sugar().Error("error in getting relative path", err)
			return err
		}
		m.DestPageDir = dir
	}
	// wrap data for base template
	data := map[string]interface{}{
		"PageTitle": baseFolderName,
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

	tmpl, execName, err := getTemplate(ctx, tmptPath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	// TODO: if there is a filemeta to change the base folder name give it more preced

	opBaseFile := filepath.Join(viper.GetString("filepath.destMDRoot"), "tags", "tags.html")
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

	tmpl, execName, err := getTemplate(ctx, tagMeta[0].TagTemplatePath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	// TODO: if there is a filemeta to change the base folder name give it more preced

	// opBaseFile := filepath.Join(viper.GetString("filepath.destMDRoot"), baseFolderName, strings.TrimSuffix(baseFolderName, filepath.Ext(baseFolderName))+".html")
	opBaseFile := baseFolderName
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
