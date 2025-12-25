package main_test

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func newTestViper(configDir string) error {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	return viper.ReadInConfig()
}

type TestEnv struct {
	Ctx context.Context
	Tmp string
}

func setTest(t *testing.T, configDir string) *TestEnv {
	t.Helper()

	ctx := context.Background()

	err := newTestViper(configDir)
	if err != nil {
		t.Fatalf("config load failed: %v", err)
	}

	logger, err := logs.InitializeLogger()
	if err != nil {
		t.Fatalf("logger init failed: %v", err)
	}

	ctx = logs.SetLoggerctx(ctx, logger)

	return &TestEnv{
		Ctx: ctx,
	}
}

// func setTest(configPath string) (context.Context, error) {
// 	ctx := context.Background()
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("json")
// 	viper.AddConfigPath("./configurables") // path to look for the config file in

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		return ctx, nil
// 	}
// 	l, err := logs.InitializeLogger()
// 	if err != nil {
// 		return ctx, err
// 	}
// 	ctx = logs.SetLoggerctx(ctx, l)
// 	ctx = logs.SetLoggerctx(ctx, l)

// 	// funcMap := template.FuncMap{
// 	// 	"config": func(key string) any {
// 	// 		return viper.Get(key)
// 	// 	},
// 	// }

// 	// _, err = template.New("").Funcs(funcMap).ParseGlob("./templates/*.html")
// 	// if err != nil {
// 	// 	return ctx, err
// 	// }
// 	return ctx, nil
// }

// func TestParseMarkdownToHtml(t *testing.T) {
// 	mdDir := "./test.md"
// 	//	destDir := "./test.html"
// 	ctx, err := setTest()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	htlmString, err := parser.ParseMarkdownToHtml(ctx, mdDir)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	fmt.Println(htlmString)
// }

func TestExtractHeadingInMarkdown(t *testing.T) {
	mdDir := "./test.md"
	ctx := context.Background()
	_, err := parser.ExtractHeadingInMarkdown(ctx, mdDir)
	if err != nil {
		t.Error(err.Error())
	}
}

// func TestSaveBaseFile(t *testing.T) {

// 	ctx, err := setTest()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = parser.ParseFiles(ctx)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// }

// func TestCopyFiles(t *testing.T) {
// 	_, err := setTest()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = util.CopyAssets(viper.GetString("filepath.mdAssetsSourcePath"), viper.GetString("filepath.mdAssetsDestPath"))
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// }

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
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := http.FileServer(http.Dir("./testdata/basic/output"))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// func setTest(t *testing.T, configDir string) *TestEnv {
// 	t.Helper()

// 	ctx := context.Background()

// 	v, err := newTestViper(configDir)
// 	if err != nil {
// 		t.Fatalf("config load failed: %v", err)
// 	}

// 	logger, err := logs.InitializeLogger()
// 	if err != nil {
// 		t.Fatalf("logger init failed: %v", err)
// 	}

// 	ctx = logs.SetLoggerctx(ctx, logger)

// 	return &TestEnv{
// 		Ctx:   ctx,
// 		Viper: v,
// 	}
// }

func TestServeGeneratedFiles(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := http.FileServer(http.Dir("/home/void/Voidzone/hanamark/test/test_output/point_B_01"))
	handler.ServeHTTP(w, req)
	// python3 -m http.server 3333

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBuildFileFlag(t *testing.T) {

	tests := []struct {
		name      string
		configDir string
	}{
		{
			name:      "basic site",
			configDir: "./test_data/01/configurables/",
		},
		{
			name:      "markdown-blog",
			configDir: "./test_data/02/configurables/",
		},

		{
			name:      "nested list basic site",
			configDir: "./test_data/03/configurables/",
		},
		{
			name:      "indexContent config - home section as index.html",
			configDir: "./test_data/04/configurables/",
		},
		{
			name:      "indexContent config - single page (about) as index.html",
			configDir: "./test_data/05/configurables/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setTest(t, tt.configDir)

			sourceStatic := viper.GetString("filepath.sourceStaticFiles")

			if sourceStatic != "" {
				err := util.CopyAssets(sourceStatic, filepath.Join(viper.GetString("filepath.destHtmlDir"), "static"))
				if err != nil {
					t.Fatalf("copy static files failed: %v", err)
				}
			}
			sourceAssets := viper.GetString("filepath.sourceAssetsPath")
			if sourceAssets != "" {
				err := util.CopyAssets(sourceAssets, filepath.Join(viper.GetString("filepath.destHtmlDir"), "static"))
				if err != nil {
					t.Fatalf("copy assets failed: %v", err)
				}
			}
			err := parser.ParseFiles(env.Ctx)
			if err != nil {
				t.Fatalf("parse files failed: %v", err)
			}
		})
	}
}
