package testing

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"log"
	"testing"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func setTest() (context.Context, error) {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Println("there is a error loading environment variables", err)
		return ctx, err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./") // path to look for the config file in

	err = viper.ReadInConfig()
	if err != nil {
		return ctx, nil
	}
	l, err := logs.InitializeLogger()
	if err != nil {
		return ctx, err
	}
	ctx = logs.SetLoggerctx(ctx, l)
	_, err = template.ParseGlob("../templates/*.html")
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func TestParseMarkdownToHtml(t *testing.T) {
	mdDir := "./test.md"
	//	destDir := "./test.html"
	htlmString, err := parser.ParseMarkdownToHtml(mdDir)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(htlmString)
}

func TestExtractHeadingInMarkdown(t *testing.T) {
	mdDir := "./test.md"
	ctx := context.Background()
	_, err := parser.ExtractHeadingInMarkdown(ctx, mdDir)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSaveBaseFile(t *testing.T) {

	ctx, err := setTest()
	if err != nil {
		t.Error(err)
	}
	err = parser.ParseFiles(ctx)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCopyFiles(t *testing.T) {
	_, err := setTest()
	if err != nil {
		t.Error(err)
	}
	err = util.CopyImages(viper.GetString("filepath.sourceImagePath"), viper.GetString("filepath.destImagePath"))
	if err != nil {
		t.Errorf(err.Error())
	}
}
