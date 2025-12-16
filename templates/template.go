package tmplt

import (
	"context"
	logs "hanamark/logger"
	"hanamark/model"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// takes in the base template and appends the content the base template and gives us back the final html string
func RenderTemplate(ctx context.Context, meta *model.PageMeta, templatePath string) error {
	l := logs.GetLoggerctx(ctx)

	// templateKey := meta.BaseFile

	// templateMap := viper.GetStringMapString("fileMeta.templateMap")
	// baseTemplatehtml, ok := templateMap[templateKey]
	// if !ok {
	// 	// there is no templating configured, so the input generated html is the output rendered template
	// 	return meta.GenHtml, nil
	// }

	// // path := filepath.Join(viper.GetString("filepath.templatePath"), baseTemplatehtml)
	// // if _, err := os.Stat(path); os.IsNotExist(err) {
	// // 	fmt.Println(err)
	// // }

	// // TODO: we got to write it to the template from templatePath
	// // Parse all templates, but only execute the ones needed
	// tmpl, err := template.ParseGlob(filepath.Join(viper.GetString("filepath.templatePath"), "*.html"))
	// if err != nil {
	// 	l.Sugar().Error("Template parsing error:", err)
	// 	return "", err
	// }

	// var buf bytes.Buffer
	// err = tmpl.ExecuteTemplate(&buf, baseTemplatehtml, meta) // i could have directly written it into the html but i am retarded
	// if err != nil {
	// 	l.Sugar().Error("Error executing template", err)
	// 	return "", err
	// }

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	// TODO: if there is a filemeta to change the base folder name give it more preced

	opFile := meta.DestPageDir
	f, err := os.Create(opFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, meta)
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return err
	}

	return nil
}

// RenderBaseLinkTemplate processes files which has links of other sub files. ususlly every folder has individual files
// and the files are linked if it is a sub folder!
func RenderBaseLinkTemplate(ctx context.Context, meta []*model.PageMeta, lp *model.ListPage) error {
	l := logs.GetLoggerctx(ctx)
	baseFolderName := lp.Base
	// templateKey := basefileName

	// templateMap := viper.GetStringMapString("fileMeta.templateMap")
	// baseTemplatehtml, ok := templateMap[templateKey]
	// if !ok {
	// 	return errors.New("base template not configured")
	// }

	// Parse all templates, but only execute the ones needed
	// tmpl, err := template.ParseGlob(filepath.Join(viper.GetString("filepath.templatePath"), "*.html"))
	// if err != nil {
	// 	l.Sugar().Error("Template parsing error:", err)
	// 	return err
	// }
	tmpl, err := template.ParseFiles(lp.TempPath)
	if err != nil {
		l.Sugar().Error("Template parsing error:", err)
		return err
	}
	// TODO: if there is a filemeta to change the base folder name give it more preced

	opBaseFile := filepath.Join(viper.GetString("filepath.destMDRoot"), baseFolderName, strings.TrimSuffix(baseFolderName, filepath.Ext(baseFolderName))+".html")
	f, err := os.Create(opBaseFile)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}
	defer f.Close()
	for _, m := range meta {
		dir, err := filepath.Rel(baseFolderName, m.DestPageDir)
		if err != nil {
			l.Sugar().Error("error in getting relative path", err)
			return err
		}
		m.DestPageDir = dir
	}
	err = tmpl.Execute(f, meta)
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return err
	}

	return nil
}

func RenderBaseTagListTemplate(ctx context.Context, taglist []*model.TagList, tmptPath string) error {
	l := logs.GetLoggerctx(ctx)

	// templateKey := basefileName

	// templateMap := viper.GetStringMapString("fileMeta.templateMap")
	// baseTemplatehtml, ok := templateMap[templateKey]
	// if !ok {
	// 	return errors.New("base template not configured")
	// }

	// Parse all templates, but only execute the ones needed
	// tmpl, err := template.ParseGlob(filepath.Join(viper.GetString("filepath.templatePath"), "*.html"))
	// if err != nil {
	// 	l.Sugar().Error("Template parsing error:", err)
	// 	return err
	// }
	tmpl, err := template.ParseFiles(tmptPath)
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

	err = tmpl.Execute(f, taglist)
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return err
	}

	return nil
}

func RenderTagLinkTemplate(ctx context.Context, tagMeta []*model.Tag, tagName string) error {
	l := logs.GetLoggerctx(ctx)
	baseFolderName := tagMeta[0].TagDestPath

	tmpl, err := template.ParseFiles(tagMeta[0].TagTemplatePath)
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

	err = tmpl.Execute(f, tagMeta)
	if err != nil {
		l.Sugar().Error("Error executing tag meta template", err)
		return err
	}

	return nil
}

func WriteIntoFile(ctx context.Context, input string, meta *model.PageMeta) error {
	l := logs.GetLoggerctx(ctx)

	f, err := os.Create(meta.DestPageDir)
	if err != nil {
		l.Sugar().Error("file creation failed", err)
		return err
	}

	defer f.Close()
	_, err = f.Write([]byte(input))
	if err != nil {
		l.Sugar().Error("writing into the file failed", err)
		return err
	}

	return nil
}
