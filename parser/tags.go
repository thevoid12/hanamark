package parser

import (
	"context"
	"errors"
	logs "hanamark/logger"
	"hanamark/model"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func ProcessTags(ctx context.Context, tag string) (tagMeta *model.Tag, err error) {
	l := logs.GetLoggerctx(ctx)

	if tag == "" {
		l.Sugar().Error("cannot parse empty tags")
		return nil, errors.New("cannot parse empty tags")
	}

	templateRootPath := viper.GetString("Filepath.templatePath")
	if templateRootPath == "" {
		return nil, errors.New("templatePath is empty in config")
	}
	tagPath := filepath.Join(templateRootPath, "tags")

	tagTemplatePath := ""
	tp := filepath.Join(tagPath, tag, "list.html")
	_, err = os.Stat(tag)
	if err == nil {
		tagTemplatePath = tp
	}
	if err != nil && os.IsNotExist(err) {
		tp = filepath.Join(tagPath, "list.html")
		_, err := os.Stat(tp) // defult list template for all the tags
		if err != nil {
			return nil, err
		}
		tagTemplatePath = tp
	}

	rootDestPath := viper.GetString("Filepath.destMDRoot")
	if rootDestPath == "" {
		return nil, errors.New("root destination path is empty in config")
	}
	destPath := filepath.Join(rootDestPath, "tags", tag+".html")
	return &model.Tag{
		TagDestPath:     destPath,
		TagName:         tag,
		TagTemplatePath: tagTemplatePath,
		FileHeading:     "", // this is filled by calculating the relative directory between the tag html and its corresponding html
		FileDestPath:    "",
	}, nil

}
