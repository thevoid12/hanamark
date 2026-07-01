package imageproc

import (
	"hanamark/model"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KarpelesLab/gowebp"
	"github.com/spf13/viper"
)

// writeMislabeledWebP writes genuine WebP content to a file with a ".png"
// extension - mirroring real-world content that got saved/exported without
// the extension matching what's actually inside (see the WebP decoder
// registration in imageproc.go's init()).
func writeMislabeledWebP(t *testing.T, path string, w, h int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create mislabeled webp: %v", err)
	}
	defer f.Close()
	if err := gowebp.Encode(f, img, &gowebp.Options{Lossy: true, Quality: 80}); err != nil {
		t.Fatalf("encode webp: %v", err)
	}
}

func writeTestPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 100, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create test png: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode test png: %v", err)
	}
}

func TestParseDestination(t *testing.T) {
	cases := []struct {
		name       string
		dest       string
		wantPath   string
		wantPreset string
		wantWidth  int
		wantHeight int
	}{
		{"no query", "./assets/hero.jpg", "./assets/hero.jpg", "", 0, 0},
		{"preset only", "./assets/hero.jpg?preset=banner", "./assets/hero.jpg", "banner", 0, 0},
		{"width only", "./assets/shot.png?w=1000", "./assets/shot.png", "", 1000, 0},
		{"height only", "./assets/portrait.png?h=1200", "./assets/portrait.png", "", 0, 1200},
		{"width wins over height when both given", "./assets/shot.png?w=1000&h=500", "./assets/shot.png", "", 1000, 0},
		{"malformed query is ignored", "./assets/x.png?%zz", "./assets/x.png", "", 0, 0},
		{"negative width ignored", "./assets/x.png?w=-5", "./assets/x.png", "", 0, 0},
		{"percent-encoded space is decoded", "../../assets/Pasted%20image%2020260527210847.png", "../../assets/Pasted image 20260527210847.png", "", 0, 0},
		{"percent-encoded path with a directive", "./assets/my%20photo.png?w=800", "./assets/my photo.png", "", 800, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path, d := ParseDestination(tc.dest)
			if path != tc.wantPath {
				t.Errorf("path = %q, want %q", path, tc.wantPath)
			}
			if d.Preset != tc.wantPreset || d.Width != tc.wantWidth || d.Height != tc.wantHeight {
				t.Errorf("directive = %+v, want preset=%q width=%d height=%d", d, tc.wantPreset, tc.wantWidth, tc.wantHeight)
			}
		})
	}
}

func TestResolveImageSource(t *testing.T) {
	tmp := t.TempDir()
	sourceMDDir := filepath.Join(tmp, "source_md")
	assetsDir := filepath.Join(sourceMDDir, "assets")
	nestedDir := filepath.Join(sourceMDDir, "blog", "2026")

	writeTestPNG(t, filepath.Join(assetsDir, "shared.png"), 4, 4)
	writeTestPNG(t, filepath.Join(nestedDir, "assets", "local.png"), 4, 4)

	viper.Reset()
	viper.Set("filepath.sourceMDDir", sourceMDDir)
	viper.Set("filepath.mdAssetsSourcePath", assetsDir)
	defer viper.Reset()

	t.Run("relative to markdown file itself", func(t *testing.T) {
		mdPath := filepath.Join(nestedDir, "post.md")
		got, ok := ResolveImageSource(mdPath, "./assets/local.png")
		if !ok {
			t.Fatal("expected resolution to succeed")
		}
		want := filepath.Join(nestedDir, "assets", "local.png")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("upward search toward sourceMDDir", func(t *testing.T) {
		mdPath := filepath.Join(nestedDir, "post.md")
		got, ok := ResolveImageSource(mdPath, "./assets/shared.png")
		if !ok {
			t.Fatal("expected resolution to succeed via upward search")
		}
		want := filepath.Join(assetsDir, "shared.png")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("mdAssetsSourcePath basename fallback", func(t *testing.T) {
		mdPath := filepath.Join(sourceMDDir, "about.md")
		got, ok := ResolveImageSource(mdPath, "somewhere/shared.png")
		if !ok {
			t.Fatal("expected resolution to succeed via mdAssetsSourcePath fallback")
		}
		want := filepath.Join(assetsDir, "shared.png")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("not found passes through", func(t *testing.T) {
		mdPath := filepath.Join(sourceMDDir, "about.md")
		_, ok := ResolveImageSource(mdPath, "./assets/does-not-exist.png")
		if ok {
			t.Error("expected ok=false for a missing file")
		}
	})

	t.Run("remote url passes through", func(t *testing.T) {
		_, ok := ResolveImageSource(filepath.Join(sourceMDDir, "about.md"), "https://example.com/x.png")
		if ok {
			t.Error("expected ok=false for a remote url")
		}
	})

	t.Run("root relative path passes through", func(t *testing.T) {
		_, ok := ResolveImageSource(filepath.Join(sourceMDDir, "about.md"), "/static/x.png")
		if ok {
			t.Error("expected ok=false for a root-relative path")
		}
	})
}

// TestParseDestinationThenResolveHandlesPercentEncodedSpaces reproduces a
// real bug found against a live site: files pasted via Obsidian's clipboard
// paste feature are named "Pasted image <timestamp>.png" (literal space on
// disk), but the markdown link percent-encodes the space as "%20". Before
// ParseDestination decoded the path, ResolveImageSource would look for a
// file literally named "...%20..." on disk, never find it, and silently
// fall back to a plain <img> tag.
func TestParseDestinationThenResolveHandlesPercentEncodedSpaces(t *testing.T) {
	tmp := t.TempDir()
	sourceMDDir := filepath.Join(tmp, "source_md")
	mdPath := filepath.Join(sourceMDDir, "blog", "2026_05", "post.md")
	realFile := filepath.Join(sourceMDDir, "assets", "Pasted image 20260527210847.png")

	writeTestPNG(t, realFile, 4, 4)

	viper.Reset()
	viper.Set("filepath.sourceMDDir", sourceMDDir)
	defer viper.Reset()

	markdownRef := "../../assets/Pasted%20image%2020260527210847.png"
	cleanPath, _ := ParseDestination(markdownRef)

	absPath, ok := ResolveImageSource(mdPath, cleanPath)
	if !ok {
		t.Fatalf("expected the percent-encoded reference %q to resolve to %q", markdownRef, realFile)
	}
	if absPath != realFile {
		t.Errorf("resolved to %q, want %q", absPath, realFile)
	}
}

func TestNormalizeBreakpoints(t *testing.T) {
	cases := []struct {
		name        string
		breakpoints []int
		target      int
		natural     int
		want        []int
	}{
		{"filters and appends target", []int{400, 800, 1200}, 800, 2000, []int{400, 800}},
		{"caps at natural width", []int{400, 800, 1200}, 1200, 600, []int{400, 600}},
		{"empty breakpoints falls back to target", nil, 800, 2000, []int{800}},
		{"dedupes when target already in list", []int{400, 800}, 800, 2000, []int{400, 800}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeBreakpoints(tc.breakpoints, tc.target, tc.natural)
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("got %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestScaledHeightPreservesAspectRatio(t *testing.T) {
	// 200x100 source (2:1) scaled to width 50 should be height 25.
	h := scaledHeight(50, 200, 100)
	if h != 25 {
		t.Errorf("scaledHeight = %d, want 25", h)
	}
}

func TestEnsureCachedSkipsSecondEncode(t *testing.T) {
	tmp := t.TempDir()
	cachePath := filepath.Join(tmp, "variant.bin")

	calls := 0
	encode := func(w io.Writer) error {
		calls++
		_, err := w.Write([]byte("payload"))
		return err
	}

	if err := ensureCached(cachePath, encode); err != nil {
		t.Fatalf("first ensureCached: %v", err)
	}
	if err := ensureCached(cachePath, encode); err != nil {
		t.Fatalf("second ensureCached: %v", err)
	}

	if calls != 1 {
		t.Errorf("encode called %d times, want 1 (second call should be a cache hit)", calls)
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("reading cached file: %v", err)
	}
	if string(data) != "payload" {
		t.Errorf("cached content = %q, want %q", data, "payload")
	}
}

func TestCacheKeyDeterminism(t *testing.T) {
	src := []byte("some source bytes")

	k1 := cacheKey(src, 800, "webp", 82)
	k2 := cacheKey(src, 800, "webp", 82)
	if k1 != k2 {
		t.Errorf("identical inputs produced different keys: %q vs %q", k1, k2)
	}

	if k3 := cacheKey(src, 400, "webp", 82); k3 == k1 {
		t.Error("different width produced the same key")
	}
	if k4 := cacheKey(src, 800, "png", 82); k4 == k1 {
		t.Error("different format produced the same key")
	}
	if k5 := cacheKey(src, 800, "webp", 60); k5 == k1 {
		t.Error("different quality produced the same key")
	}
	if k6 := cacheKey([]byte("different bytes"), 800, "webp", 82); k6 == k1 {
		t.Error("different source bytes produced the same key")
	}
}

func TestProcessImageEndToEnd(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "source_md", "assets", "hero.png")
	writeTestPNG(t, srcPath, 200, 100) // 2:1 aspect ratio

	destHtmlDir := filepath.Join(tmp, "output_html")
	mdAssetsDestPath := filepath.Join(destHtmlDir, "assets")

	viper.Reset()
	viper.Set("filepath.destHtmlDir", destHtmlDir)
	defer viper.Reset()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWD)

	result, err := ProcessImage(t.Context(), srcPath, 50, 0, []int{25, 50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	if result.Width != 50 || result.Height != 25 {
		t.Errorf("intrinsic dims = %dx%d, want 50x25", result.Width, result.Height)
	}
	if len(result.WebP) != 2 {
		t.Fatalf("expected 2 webp variants, got %d", len(result.WebP))
	}
	if len(result.Fallback) != 2 {
		t.Fatalf("expected 2 fallback variants, got %d", len(result.Fallback))
	}
	if result.FallbackExt != "png" {
		t.Errorf("fallback ext = %q, want png", result.FallbackExt)
	}
	for _, v := range append(append([]model.Variant{}, result.WebP...), result.Fallback...) {
		if !strings.HasPrefix(v.URL, "/assets/generated/") {
			t.Errorf("variant URL %q should be under /assets/generated/", v.URL)
		}
		if _, err := os.Stat(filepath.Join(destHtmlDir, strings.TrimPrefix(v.URL, "/"))); err != nil {
			t.Errorf("published variant not found on disk: %v", err)
		}
	}

	if _, err := os.Stat(filepath.Join(tmp, "hanamark_internal", "README.md")); err != nil {
		t.Errorf("expected hanamark_internal/README.md to be created: %v", err)
	}
}

// TestProcessImageDecodesMislabeledWebP guards against real content found in
// the wild: a file saved with a ".png" extension whose actual bytes are
// WebP. Before registering gowebp's decoder in init(), ProcessImage failed
// on these with "decoding source image: image: unknown format".
func TestProcessImageDecodesMislabeledWebP(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "source_md", "assets", "mislabeled.png")
	writeMislabeledWebP(t, srcPath, 100, 50) // 2:1 aspect ratio, real webp bytes

	destHtmlDir := filepath.Join(tmp, "output_html")
	mdAssetsDestPath := filepath.Join(destHtmlDir, "assets")

	viper.Reset()
	viper.Set("filepath.destHtmlDir", destHtmlDir)
	defer viper.Reset()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWD)

	result, err := ProcessImage(t.Context(), srcPath, 50, 0, []int{25, 50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("ProcessImage failed on a mislabeled webp file: %v", err)
	}
	if result.Width != 50 || result.Height != 25 {
		t.Errorf("intrinsic dims = %dx%d, want 50x25", result.Width, result.Height)
	}
}

func TestOutputFormatAndBackupOriginalFormatDefaults(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	if got := OutputFormat(); got != "webp" {
		t.Errorf("OutputFormat() with no config set = %q, want %q", got, "webp")
	}
	if !BackupOriginalFormat() {
		t.Error("BackupOriginalFormat() with no config set should default to true")
	}

	viper.Set("image.backupOriginalFormat", false)
	if BackupOriginalFormat() {
		t.Error("BackupOriginalFormat() should be false once image.backupOriginalFormat is set to false")
	}
}

func TestProcessImageSkipsFallbackWhenBackupDisabled(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "source_md", "assets", "hero.png")
	writeTestPNG(t, srcPath, 200, 100)

	destHtmlDir := filepath.Join(tmp, "output_html")
	mdAssetsDestPath := filepath.Join(destHtmlDir, "assets")

	viper.Reset()
	viper.Set("filepath.destHtmlDir", destHtmlDir)
	viper.Set("image.backupOriginalFormat", false)
	defer viper.Reset()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWD)

	result, err := ProcessImage(t.Context(), srcPath, 50, 0, []int{25, 50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	if len(result.WebP) != 2 {
		t.Fatalf("expected 2 webp variants, got %d", len(result.WebP))
	}
	if len(result.Fallback) != 0 {
		t.Errorf("expected no fallback variants when backupOriginalFormat=false, got %d", len(result.Fallback))
	}

	entries, err := os.ReadDir(filepath.Join(tmp, "hanamark_internal", "images"))
	if err != nil {
		t.Fatalf("reading cache dir: %v", err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".png") {
			t.Errorf("expected no cached .png fallback file, found %s", e.Name())
		}
	}
}

func TestUnsupportedOutputFormatFallsBackToWebP(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "source_md", "assets", "hero.png")
	writeTestPNG(t, srcPath, 100, 100)

	destHtmlDir := filepath.Join(tmp, "output_html")
	mdAssetsDestPath := filepath.Join(destHtmlDir, "assets")

	viper.Reset()
	viper.Set("filepath.destHtmlDir", destHtmlDir)
	viper.Set("image.outputFormat", "avif")
	defer viper.Reset()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWD)

	result, err := ProcessImage(t.Context(), srcPath, 50, 0, []int{50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	if len(result.WebP) != 1 {
		t.Fatalf("expected an unsupported outputFormat to fall back to webp, got %d webp variants", len(result.WebP))
	}
	if !strings.HasSuffix(result.WebP[0].URL, ".webp") {
		t.Errorf("expected the fallback-to-webp variant URL to end in .webp, got %q", result.WebP[0].URL)
	}
}

// TestPublishOverwritesInPlaceWhenSourceContentChanges guards the whole
// reason the published filename is path-derived rather than content-hash
// derived: replacing a source image's bytes (same name, same preset) must
// overwrite the same published path, not accumulate a second file alongside
// the stale one.
func TestPublishOverwritesInPlaceWhenSourceContentChanges(t *testing.T) {
	tmp := t.TempDir()
	sourceMDDir := filepath.Join(tmp, "source_md")
	srcPath := filepath.Join(sourceMDDir, "blog", "post1", "assets", "hero.png")

	destHtmlDir := filepath.Join(tmp, "output_html")
	mdAssetsDestPath := filepath.Join(destHtmlDir, "assets")

	viper.Reset()
	viper.Set("filepath.sourceMDDir", sourceMDDir)
	viper.Set("filepath.destHtmlDir", destHtmlDir)
	defer viper.Reset()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWD)

	writeTestPNG(t, srcPath, 200, 100)
	result1, err := ProcessImage(t.Context(), srcPath, 50, 0, []int{50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("first ProcessImage failed: %v", err)
	}
	if len(result1.WebP) != 1 {
		t.Fatalf("expected 1 webp variant, got %d", len(result1.WebP))
	}
	firstURL := result1.WebP[0].URL
	firstPublishedPath := filepath.Join(destHtmlDir, strings.TrimPrefix(firstURL, "/"))
	firstBytes, err := os.ReadFile(firstPublishedPath)
	if err != nil {
		t.Fatalf("reading first published file: %v", err)
	}

	// Replace the source image with different content, same name, same
	// preset/width - simulates a user swapping out a photo.
	writeTestPNG(t, srcPath, 300, 150)
	result2, err := ProcessImage(t.Context(), srcPath, 50, 0, []int{50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("second ProcessImage failed: %v", err)
	}
	if len(result2.WebP) != 1 {
		t.Fatalf("expected 1 webp variant, got %d", len(result2.WebP))
	}
	secondURL := result2.WebP[0].URL

	if secondURL != firstURL {
		t.Errorf("published URL changed after replacing source content (%q -> %q); expected the same deterministic path to be reused", firstURL, secondURL)
	}

	secondBytes, err := os.ReadFile(firstPublishedPath)
	if err != nil {
		t.Fatalf("reading published file after second run: %v", err)
	}
	if string(firstBytes) == string(secondBytes) {
		t.Error("published file content did not change after replacing the source image - overwrite did not happen")
	}

	// Default config also publishes a same-format (.png) fallback alongside
	// the .webp variant - so 2 files total per generation, still not 4 after
	// regenerating twice, which is the actual "no accumulation" assertion.
	publishedDir := filepath.Dir(firstPublishedPath)
	entries, err := os.ReadDir(publishedDir)
	if err != nil {
		t.Fatalf("reading published dir: %v", err)
	}
	if len(entries) != 2 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected exactly 2 published files (webp + fallback) after regenerating twice (no accumulation), got %d: %v", len(entries), names)
	}
}

// TestPublishDoesNotCollideAcrossSourceFoldersWithSameBasename guards the
// risk flagged when switching away from content-hash naming: source
// basenames are not unique across a whole site (e.g. two different blog
// posts each with their own "assets/001.png"), so the published path must be
// derived from the full relative source path, not just the basename.
func TestPublishDoesNotCollideAcrossSourceFoldersWithSameBasename(t *testing.T) {
	tmp := t.TempDir()
	sourceMDDir := filepath.Join(tmp, "source_md")
	src1 := filepath.Join(sourceMDDir, "blog", "post1", "assets", "001.png")
	src2 := filepath.Join(sourceMDDir, "blog", "post2", "assets", "001.png")

	destHtmlDir := filepath.Join(tmp, "output_html")
	mdAssetsDestPath := filepath.Join(destHtmlDir, "assets")

	viper.Reset()
	viper.Set("filepath.sourceMDDir", sourceMDDir)
	viper.Set("filepath.destHtmlDir", destHtmlDir)
	defer viper.Reset()

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(origWD)

	writeTestPNG(t, src1, 200, 100)
	writeTestPNG(t, src2, 300, 150) // different content, same basename "001.png"

	result1, err := ProcessImage(t.Context(), src1, 50, 0, []int{50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("ProcessImage(src1) failed: %v", err)
	}
	result2, err := ProcessImage(t.Context(), src2, 50, 0, []int{50}, mdAssetsDestPath)
	if err != nil {
		t.Fatalf("ProcessImage(src2) failed: %v", err)
	}

	url1, url2 := result1.WebP[0].URL, result2.WebP[0].URL
	if url1 == url2 {
		t.Fatalf("two different source images with the same basename published to the same URL: %q", url1)
	}

	path1 := filepath.Join(destHtmlDir, strings.TrimPrefix(url1, "/"))
	path2 := filepath.Join(destHtmlDir, strings.TrimPrefix(url2, "/"))
	bytes1, err := os.ReadFile(path1)
	if err != nil {
		t.Fatalf("reading published file 1: %v", err)
	}
	bytes2, err := os.ReadFile(path2)
	if err != nil {
		t.Fatalf("reading published file 2: %v", err)
	}
	if string(bytes1) == string(bytes2) {
		t.Error("expected different content for two different source images, got identical bytes")
	}
}
