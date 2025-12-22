package main_test

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"log"
	"net/http"
	"testing"

	"github.com/spf13/viper"
)

func setTest() (context.Context, error) {
	ctx := context.Background()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configurables") // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		return ctx, nil
	}
	l, err := logs.InitializeLogger()
	if err != nil {
		return ctx, err
	}
	ctx = logs.SetLoggerctx(ctx, l)
	ctx = logs.SetLoggerctx(ctx, l)

	// funcMap := template.FuncMap{
	// 	"config": func(key string) any {
	// 		return viper.Get(key)
	// 	},
	// }

	// _, err = template.New("").Funcs(funcMap).ParseGlob("./templates/*.html")
	// if err != nil {
	// 	return ctx, err
	// }
	return ctx, nil
}

func TestParseMarkdownToHtml(t *testing.T) {
	mdDir := "./test.md"
	//	destDir := "./test.html"
	ctx, err := setTest()
	if err != nil {
		t.Error(err)
	}
	htlmString, err := parser.ParseMarkdownToHtml(ctx, mdDir)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(htlmString)
}

func TestExtractHeadingInMarkdown(t *testing.T) {
	mdDir := "./test.md"
	ctx := context.Background()
	_, err := parser.ExtractHeadingInMarkdown(ctx, mdDir)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSaveBaseFile(t *testing.T) {

	ctx, err := setTest()
	if err != nil {
		t.Error(err)
	}
	err = parser.ParseFiles(ctx)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCopyFiles(t *testing.T) {
	_, err := setTest()
	if err != nil {
		t.Error(err)
	}
	err = util.CopyAssets(viper.GetString("filepath.sourceAssetsPath"), viper.GetString("filepath.destAssetsPath"))
	if err != nil {
		t.Error(err.Error())
	}
}

func TestParseFrontMatter(t *testing.T) {
	mdDir := "./test.md"
	ctx := context.Background()
	fm, err := parser.ParseFrontMatter(ctx, mdDir)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(fm)
}

func TestServeStaticFiles(t *testing.T) {

	// ctx, err := setTest()
	// if err != nil {
	// 	t.Error(err)
	// }
	serveStaticFiles("./point_B")
	//	if err != nil {
	//		t.Error(err.Error())
	//	}
}

func serveStaticFiles(dir string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
