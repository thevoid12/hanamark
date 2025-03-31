package tmplt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"
)

// takes in the base template and appends the content the base template and gives us back the final html string
func RenderTemplate(ctx context.Context, meta *model.PageMeta) (string, error) {
	l := logs.GetLoggerctx(ctx)

	templateKey := meta.BaseFile

	templateMap := viper.GetStringMapString("fileMeta.templateMap")
	baseTemplatehtml, ok := templateMap[templateKey]
	if !ok {
		// there is no templating configured, so the input generated html is the output rendered template
		return meta.GenHtml, nil
	}

	path := filepath.Join(viper.GetString("filepath.templatePath"), baseTemplatehtml)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		l.Sugar().Error("this type of file is not configured in config so template cannot be rendered", err)
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, meta) // i could have directly written into the html but i am retarded
	if err != nil {
		l.Sugar().Error("Error executing template", err)
		return "", err
	}

	return buf.String(), nil
}

func RenderBaseTemplate(ctx context.Context, meta []*model.PageMeta, basefileName string) error {
	l := logs.GetLoggerctx(ctx)

	templateKey := basefileName

	templateMap := viper.GetStringMapString("fileMeta.templateMap")
	baseTemplatehtml, ok := templateMap[templateKey]
	if !ok {
		return errors.New("base template not configured")
	}

	baseTemplatepath := filepath.Join(viper.GetString("filepath.templatePath"), baseTemplatehtml)

	tmpl, err := template.ParseFiles(baseTemplatepath)
	if err != nil {
		l.Sugar().Error("this type of file is not configured in config so template cannot be rendered", err)
		return err
	}
	opBaseFile := filepath.Join(viper.GetString("filepath.destMDRoot"), basefileName)
	f, err := os.Create(opBaseFile)
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
