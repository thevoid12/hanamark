package tmplt

import (
	"bytes"
	"context"
	logs "hanamark/logger"
	"hanamark/model"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

// takes in the base template and appends the content the base template and gives us back the final html string
func RenderTemplate(ctx context.Context, meta *model.PageMeta) (string, error) {
	l := logs.GetLoggerctx(ctx)

	templateMap := viper.GetStringMapString("filepath.templateMap")
	baseTemplatehtml, ok := templateMap[meta.PageType]
	if !ok {
		// there is no templating configured, so the input generated html is the output rendered template
		return meta.GenHtml, nil
	}

	tmpl, err := template.ParseFiles(baseTemplatehtml)
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
