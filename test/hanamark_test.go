package main_test

import (
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
		{
			name:      "indexContent config - nested folder1 section as index.html",
			configDir: "./test_data/06/configurables/",
		},
		{
			name:      "opengraph meta tags feature",
			configDir: "./test_data/07/configurables/",
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

// TestImageProcessingPipeline exercises the image-optimization pipeline end
// to end against test_data/08: an un-annotated page (both images default to
// the "content" preset), a page opting its first image into "banner" via
// frontmatter, a page using per-image markdown query-string overrides, and
// the image.enabled=false fallback which must reproduce today's plain <img>
// output byte-for-byte.
func TestImageProcessingPipeline(t *testing.T) {
	configDir := "./test_data/08/configurables/"
	env := setTest(t, configDir)

	if err := parser.ParseFiles(env.Ctx); err != nil {
		t.Fatalf("parse files failed: %v", err)
	}

	destHtmlDir := viper.GetString("filepath.destHtmlDir")
	readOutput := func(t *testing.T, name string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(destHtmlDir, name))
		if err != nil {
			t.Fatalf("reading generated %s: %v", name, err)
		}
		return string(data)
	}

	t.Run("un-annotated images use the content preset", func(t *testing.T) {
		html := readOutput(t, "about.html")
		if strings.Count(html, "<picture>") != 2 {
			t.Fatalf("expected 2 <picture> blocks, got:\n%s", html)
		}
		if strings.Contains(html, `fetchpriority="high"`) {
			t.Error("no image on this page opted into banner treatment, but fetchpriority=\"high\" was found")
		}
		if strings.Count(html, `loading="lazy"`) != 2 {
			t.Errorf("expected both images to be loading=\"lazy\" (content preset):\n%s", html)
		}
	})

	t.Run("first_image_preset banner opts the first image in", func(t *testing.T) {
		html := readOutput(t, "banner.html")
		firstIdx := strings.Index(html, "<picture>")
		secondIdx := strings.Index(html[firstIdx+1:], "<picture>")
		if firstIdx == -1 || secondIdx == -1 {
			t.Fatalf("expected 2 <picture> blocks, got:\n%s", html)
		}
		firstBlock := html[:firstIdx+secondIdx]
		if !strings.Contains(firstBlock, `loading="eager"`) || !strings.Contains(firstBlock, `fetchpriority="high"`) {
			t.Errorf("expected first image to use the banner preset (eager/high), got:\n%s", firstBlock)
		}
		secondBlock := html[firstIdx+secondIdx:]
		if !strings.Contains(secondBlock, `loading="lazy"`) || strings.Contains(secondBlock, `fetchpriority`) {
			t.Errorf("expected second image to use the content preset (lazy, no fetchpriority), got:\n%s", secondBlock)
		}
	})

	t.Run("markdown directive overrides win", func(t *testing.T) {
		html := readOutput(t, "override.html")
		if !strings.Contains(html, `fetchpriority="high"`) {
			t.Errorf("expected ?preset=banner to apply the banner preset, got:\n%s", html)
		}
		if !strings.Contains(html, `width="90"`) {
			t.Errorf("expected ?w=90 to set an explicit width, got:\n%s", html)
		}
	})

	t.Run("image.enabled=false falls back to plain img", func(t *testing.T) {
		viper.Set("image.enabled", false)
		defer viper.Set("image.enabled", true)

		if err := parser.ParseFiles(env.Ctx); err != nil {
			t.Fatalf("parse files failed: %v", err)
		}

		html := readOutput(t, "about.html")
		if strings.Contains(html, "<picture>") {
			t.Errorf("expected no <picture> markup when image.enabled=false, got:\n%s", html)
		}
		if !strings.Contains(html, `<img src="./assets/photo1.png" alt="First photo" />`) {
			t.Errorf("expected plain passthrough <img> tag, got:\n%s", html)
		}
	})
}
