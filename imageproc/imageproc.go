// Package imageproc converts markdown content images into resized, cached
// WebP (plus a same-format raster fallback) picture markup at build time, so
// pages ship responsive images with explicit width/height (fixing CLS)
// without the generator ever needing to know the visitor's screen size.
package imageproc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"hanamark/util"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/KarpelesLab/gowebp"
	"github.com/gomarkdown/markdown/ast"
	gmhtml "github.com/gomarkdown/markdown/html"
	"github.com/spf13/viper"
	"golang.org/x/image/draw"
)

const (
	// PresetBanner and PresetContent are the only recognized preset names.
	// A page's first image only uses PresetBanner when its frontmatter opts
	// in via first_image_preset: banner; every other image uses PresetContent.
	PresetBanner  = "banner"
	PresetContent = "content"

	internalDir        = "hanamark_internal"
	internalImagesSub  = "images"
	defaultQuality     = 82
	generatedURLSubdir = "generated"
)

const internalReadme = `This directory is a build cache managed automatically by hanamark.

Do not delete it if you want fast rebuilds -- resizing and re-encoding images
is the slowest part of a build, and this cache lets unchanged images skip
that work on every subsequent "hanamark build".

That said, it is completely safe to delete. If this folder is missing (or
partially deleted), hanamark will transparently and automatically recreate
whatever is needed on the next build. There is nothing to restore or repair
by hand.
`

func init() {
	// gowebp already provides both the encoder we use and a decoder; register
	// it with the stdlib image package so a source image that is actually
	// WebP content (regardless of its file extension - e.g. a ".png" that's
	// really WebP) decodes correctly instead of erroring as "unknown format".
	image.RegisterFormat("webp", "RIFF????WEBP", gowebp.Decode, gowebp.DecodeConfig)
}

// Enabled reports whether image processing is turned on. Defaults to true
// when the "image.enabled" key is absent from config.json.
func Enabled() bool {
	if !viper.IsSet("image.enabled") {
		return true
	}
	return viper.GetBool("image.enabled")
}

// PresetExists reports whether image.presets.<name> is defined in config.
func PresetExists(name string) bool {
	return viper.IsSet("image.presets." + name)
}

func presetKey(name, field string) string {
	return "image.presets." + name + "." + field
}

// PresetWidth returns the configured width for a preset, or 0 if unset.
func PresetWidth(name string) int {
	return viper.GetInt(presetKey(name, "width"))
}

// PresetBreakpoints returns the configured srcset breakpoint widths for a preset.
func PresetBreakpoints(name string) []int {
	return viper.GetIntSlice(presetKey(name, "breakpoints"))
}

// PresetLoading returns the loading="" attribute value for a preset (eager/lazy).
func PresetLoading(name string) string {
	return viper.GetString(presetKey(name, "loading"))
}

// PresetFetchPriority returns the fetchpriority="" attribute value for a preset.
// "auto" (or empty) means the attribute is omitted entirely.
func PresetFetchPriority(name string) string {
	return viper.GetString(presetKey(name, "fetchpriority"))
}

// isValidFetchPriority reports whether v is one of the fetchpriority values
// defined by the HTML spec. Anything else (including empty, meaning no
// override was given) is rejected so a malformed ?fetchpriority= directive
// silently falls back to the preset's own value instead of emitting garbage
// into the generated markup.
func isValidFetchPriority(v string) bool {
	switch v {
	case "high", "low", "auto":
		return true
	default:
		return false
	}
}

// PresetSizes returns the sizes="" attribute value for a preset.
func PresetSizes(name string) string {
	return viper.GetString(presetKey(name, "sizes"))
}

// Quality returns the configured encode quality for a format, defaulting to 82.
func Quality(format string) int {
	q := viper.GetInt("image.quality." + format)
	if q <= 0 {
		return defaultQuality
	}
	return q
}

// OutputFormat returns the single modern format every image is converted to
// (only "webp" is supported today; "avif" is a planned future addition, not
// a set of simultaneously-generated formats - a build only ever produces one
// modern-format variant per source image, hence a single value rather than
// a list).
func OutputFormat() string {
	f := strings.ToLower(viper.GetString("image.outputFormat"))
	if f == "" {
		return "webp"
	}
	return f
}

// BackupOriginalFormat reports whether a same-format-as-source raster
// fallback (jpg/png) is generated alongside the modern-format output, for
// browsers that don't support it. Defaults to true (today's behavior) when
// "image.backupOriginalFormat" is absent from config.json. Turning it off
// avoids the extra resize/encode work and the extra files on disk for sites
// that don't need a legacy fallback.
func BackupOriginalFormat() bool {
	if !viper.IsSet("image.backupOriginalFormat") {
		return true
	}
	return viper.GetBool("image.backupOriginalFormat")
}

// ParseDestination splits a markdown image destination into the clean asset
// path (query string stripped, percent-encoding decoded) and the parsed
// model.Directive. A destination with no query string returns a zero
// model.Directive.
//
// Percent-decoding matters because authoring tools (e.g. Obsidian's
// paste-image-from-clipboard, which names files "Pasted image <timestamp>.png")
// routinely percent-encode spaces and other special characters in the
// markdown link even though the file on disk keeps its literal name - without
// decoding, the filesystem lookup in ResolveImageSource never finds the file
// and silently falls back to a plain <img> tag.
func ParseDestination(dest string) (string, model.Directive) {
	qIdx := strings.IndexByte(dest, '?')
	path := dest
	rawQuery := ""
	if qIdx != -1 {
		path = dest[:qIdx]
		rawQuery = dest[qIdx+1:]
	}
	if decoded, err := url.PathUnescape(path); err == nil {
		path = decoded
	}

	if qIdx == -1 {
		return path, model.Directive{}
	}
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return path, model.Directive{}
	}

	d := model.Directive{Preset: values.Get("preset")}
	if w := values.Get("w"); w != "" {
		if n, err := strconv.Atoi(w); err == nil && n > 0 {
			d.Width = n
		}
	}
	if h := values.Get("h"); h != "" {
		if n, err := strconv.Atoi(h); err == nil && n > 0 {
			d.Height = n
		}
	}
	if fp := values.Get("fetchpriority"); isValidFetchPriority(fp) {
		d.FetchPriority = fp
	}
	if d.Width > 0 && d.Height > 0 {
		// Resize-only (no crop) can't honor both without distorting the
		// aspect ratio - width wins, height is dropped.
		d.Height = 0
	}
	return path, d
}

// IsLocalImagePath reports whether rel is a path this pipeline is meant to
// process at all - false for empty strings, remote URLs, and root-relative
// paths (e.g. "/static/..."), which are intentionally passed through
// untouched and shouldn't be logged as resolution failures.
func IsLocalImagePath(rel string) bool {
	return rel != "" && !strings.Contains(rel, "://") && !strings.HasPrefix(rel, "/")
}

// ResolveImageSource turns a markdown image path into an absolute file on
// disk. It returns ok=false for anything that should pass through untouched:
// remote URLs, root-relative paths (e.g. "/static/..."), or files that can't
// be found anywhere.
func ResolveImageSource(mdSourcePath, rel string) (absPath string, ok bool) {
	if !IsLocalImagePath(rel) {
		return "", false
	}

	mdDir := filepath.Dir(mdSourcePath)
	if p := filepath.Join(mdDir, rel); fileExists(p) {
		return p, true
	}

	// Search upward toward sourceMDDir, matching the convention (seen
	// throughout the existing source_md tree) of every page writing
	// "./assets/x.png" against one shared assets folder regardless of
	// how deeply the page itself is nested.
	sourceMDDir := filepath.Clean(viper.GetString("filepath.sourceMDDir"))
	for dir := filepath.Dir(mdDir); ; dir = filepath.Dir(dir) {
		if p := filepath.Join(dir, rel); fileExists(p) {
			return p, true
		}
		if dir == sourceMDDir {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root without finding sourceMDDir; stop
		}
	}

	if assetsDir := viper.GetString("filepath.mdAssetsSourcePath"); assetsDir != "" {
		if p := filepath.Join(assetsDir, filepath.Base(rel)); fileExists(p) {
			return p, true
		}
	}

	return "", false
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

// publicRelPath computes the identity used for a source image's *published*
// filename: its path relative to sourceMDDir, rather than a content hash.
// Source basenames alone aren't unique across a whole site (e.g. two
// different blog posts can each have their own "assets/001.png"), so the
// relative path is what keeps two different images from colliding once
// published under a single flat mdAssetsDestPath/generated/ tree. Falls back
// to just the basename if srcAbsPath isn't actually under sourceMDDir (an
// unusual config), which trades back some collision safety in that rare
// case only.
func publicRelPath(srcAbsPath string) string {
	sourceMDDir := filepath.Clean(viper.GetString("filepath.sourceMDDir"))
	rel, err := filepath.Rel(sourceMDDir, srcAbsPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return filepath.Base(srcAbsPath)
	}
	return rel
}

// ProcessImage decodes srcAbsPath, resizes it to the breakpoint widths
// implied by targetWidth/targetHeight (whichever is > 0; the other dimension
// is always derived from the source's aspect ratio - there is no cropping),
// encodes each size to WebP plus the source's original raster format, and
// caches the results under hanamark_internal/images so unchanged images
// aren't reprocessed on subsequent builds. Generated files are published
// into mdAssetsDestPath/generated/.
func ProcessImage(ctx context.Context, srcAbsPath string, targetWidth, targetHeight int, breakpoints []int, mdAssetsDestPath string) (*model.Result, error) {
	l := logs.GetLoggerctx(ctx)

	srcBytes, err := os.ReadFile(srcAbsPath)
	if err != nil {
		return nil, fmt.Errorf("reading source image: %w", err)
	}

	srcImg, format, err := image.Decode(bytes.NewReader(srcBytes))
	if err != nil {
		return nil, fmt.Errorf("decoding source image: %w", err)
	}
	bounds := srcImg.Bounds()
	naturalW, naturalH := bounds.Dx(), bounds.Dy()
	if naturalW == 0 || naturalH == 0 {
		return nil, fmt.Errorf("source image %s has a zero dimension", srcAbsPath)
	}

	effectiveWidth := targetWidth
	if effectiveWidth <= 0 && targetHeight > 0 {
		effectiveWidth = int(math.Round(float64(targetHeight) * float64(naturalW) / float64(naturalH)))
	}
	if effectiveWidth <= 0 {
		effectiveWidth = naturalW
	}

	widths := normalizeBreakpoints(breakpoints, effectiveWidth, naturalW)
	if err := ensureInternalDir(); err != nil {
		return nil, fmt.Errorf("preparing image cache: %w", err)
	}

	outputFormat := OutputFormat()
	if outputFormat != "webp" {
		if l != nil {
			l.Sugar().Warnf("image.outputFormat %q is not supported yet, falling back to \"webp\"", outputFormat)
		}
		outputFormat = "webp"
	}
	backupOriginal := BackupOriginalFormat()

	fallbackExt := rasterExt(format)
	outputQuality := Quality(outputFormat)
	slug := strings.TrimSuffix(filepath.Base(srcAbsPath), filepath.Ext(srcAbsPath))
	publicRel := publicRelPath(srcAbsPath)

	result := &model.Result{FallbackExt: fallbackExt}
	for i, w := range widths {
		h := scaledHeight(w, naturalW, naturalH)
		scaled := srcImg
		if w != naturalW {
			scaled = resize(srcImg, w, h)
		}

		if variant, err := encodeAndPublish(scaled, srcBytes, slug, publicRel, w, outputFormat, outputQuality, mdAssetsDestPath, func(out io.Writer, img image.Image) error {
			return gowebp.Encode(out, img, &gowebp.Options{Lossy: true, Quality: float32(outputQuality)})
		}); err != nil {
			if l != nil {
				l.Sugar().Warnf("%s encode failed for %s at %dw: %v", outputFormat, srcAbsPath, w, err)
			}
		} else {
			result.WebP = append(result.WebP, variant)
		}

		if backupOriginal {
			if variant, err := encodeAndPublish(scaled, srcBytes, slug, publicRel, w, fallbackExt, outputQuality, mdAssetsDestPath, func(out io.Writer, img image.Image) error {
				return encodeRaster(out, img, fallbackExt, outputQuality)
			}); err != nil {
				if l != nil {
					l.Sugar().Warnf("fallback encode failed for %s at %dw: %v", srcAbsPath, w, err)
				}
			} else {
				result.Fallback = append(result.Fallback, variant)
			}
		}

		if i == len(widths)-1 {
			result.Width, result.Height = w, h
		}
	}

	if len(result.WebP) == 0 && len(result.Fallback) == 0 {
		return nil, fmt.Errorf("no variants could be produced for %s", srcAbsPath)
	}
	return result, nil
}

func normalizeBreakpoints(breakpoints []int, targetWidth, naturalWidth int) []int {
	maxWidth := targetWidth
	if naturalWidth < maxWidth {
		maxWidth = naturalWidth
	}
	if maxWidth <= 0 {
		maxWidth = naturalWidth
	}

	seen := map[int]bool{}
	var widths []int
	for _, bp := range breakpoints {
		if bp <= 0 || bp > maxWidth || seen[bp] {
			continue
		}
		seen[bp] = true
		widths = append(widths, bp)
	}
	if !seen[maxWidth] {
		widths = append(widths, maxWidth)
	}
	sort.Ints(widths)
	return widths
}

func scaledHeight(w, naturalW, naturalH int) int {
	if naturalW == 0 {
		return 0
	}
	return int(math.Round(float64(w) * float64(naturalH) / float64(naturalW)))
}

func resize(src image.Image, w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func rasterExt(format string) string {
	switch format {
	case "jpeg":
		return "jpg"
	default:
		// png is lossless and safe as a default fallback for anything
		// that isn't jpeg (including gif, which loses animation here -
		// static image processing only, matching the resize-only scope).
		return "png"
	}
}

func encodeRaster(w io.Writer, img image.Image, ext string, quality int) error {
	if ext == "jpg" {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	}
	return png.Encode(w, img)
}

func cacheKey(srcBytes []byte, width int, format string, quality int) string {
	h := sha256.New()
	h.Write(srcBytes)
	fmt.Fprintf(h, "|%d|%s|%d", width, format, quality)
	return hex.EncodeToString(h.Sum(nil))[:12]
}

func variantFilename(slug, hash string, width int, ext string) string {
	return fmt.Sprintf("%s-%s-%dw.%s", slug, hash, width, ext)
}

func internalImagesDir() string {
	return filepath.Join(internalDir, internalImagesSub)
}

func ensureInternalDir() error {
	if err := os.MkdirAll(internalImagesDir(), 0o755); err != nil {
		return err
	}
	readmePath := filepath.Join(internalDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		_ = os.WriteFile(readmePath, []byte(internalReadme), 0o644)
	}
	return nil
}

// ensureCached writes encode's output to cachePath only if it doesn't
// already exist, via a temp-file-then-rename so a build interrupted
// mid-write can never leave a corrupt file behind that a later build would
// mistake for a valid cache hit.
func ensureCached(cachePath string, encode func(io.Writer) error) error {
	if fileExists(cachePath) {
		return nil
	}
	tmp := cachePath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := encode(f); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, cachePath)
}

func encodeAndPublish(img image.Image, srcBytes []byte, slug, publicRel string, width int, ext string, quality int, mdAssetsDestPath string, encode func(io.Writer, image.Image) error) (model.Variant, error) {
	// Cache filenames stay content-hash-keyed so unchanged images skip
	// re-encoding; this cache lives under hanamark_internal/, is never
	// shipped, and isn't worth optimizing for accumulation.
	hash := cacheKey(srcBytes, width, ext, quality)
	cacheFilename := variantFilename(slug, hash, width, ext)
	cachePath := filepath.Join(internalImagesDir(), cacheFilename)

	if err := ensureCached(cachePath, func(w io.Writer) error { return encode(w, img) }); err != nil {
		return model.Variant{}, err
	}

	publicURL, err := publish(cachePath, publicRel, width, ext, mdAssetsDestPath)
	if err != nil {
		return model.Variant{}, err
	}
	return model.Variant{Width: width, URL: publicURL}, nil
}

// publish copies a cached variant into mdAssetsDestPath/generated/, mirroring
// the source image's own relative path (see publicRelPath) instead of a
// content hash. Regenerating the same source image at the same width always
// overwrites the same destination path rather than accumulating a new one
// alongside the old, which is what keeps this directory from growing
// unbounded as images get replaced over the site's lifetime.
func publish(cachePath, publicRel string, width int, ext string, mdAssetsDestPath string) (string, error) {
	relDir := filepath.Dir(publicRel)
	baseSlug := strings.TrimSuffix(filepath.Base(publicRel), filepath.Ext(publicRel))
	filename := fmt.Sprintf("%s-%dw.%s", baseSlug, width, ext)

	destDir := filepath.Join(mdAssetsDestPath, generatedURLSubdir, relDir)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}
	destPath := filepath.Join(destDir, filename)
	if err := util.CopyFile(cachePath, destPath); err != nil {
		return "", err
	}

	destHtmlDir := filepath.Clean(viper.GetString("filepath.destHtmlDir"))
	rel, err := filepath.Rel(destHtmlDir, destPath)
	if err != nil {
		rel = filepath.Join("assets", generatedURLSubdir, relDir, filename)
	}
	return "/" + filepath.ToSlash(rel), nil
}

// imageHookState is per-page render state: gomarkdown calls ParseMarkdownToHtml
// once per markdown file, so a fresh instance (and its seenFirst counter) is
// created for every page.
type imageHookState struct {
	ctx              context.Context
	mdSourcePath     string
	mdAssetsDestPath string
	wantsBannerFirst bool
	seenFirst        bool
}

// NewImageRenderHook returns a gomarkdown html.RenderNodeFunc that intercepts
// *ast.Image nodes and replaces gomarkdown's plain <img> output with a
// responsive <picture> block; every other node type is left untouched by
// returning (ast.GoToNext, false) so gomarkdown renders it exactly as it
// always has.
func NewImageRenderHook(ctx context.Context, mdSourcePath, mdAssetsDestPath string, wantsBannerForFirst bool) func(io.Writer, ast.Node, bool) (ast.WalkStatus, bool) {
	st := &imageHookState{
		ctx:              ctx,
		mdSourcePath:     mdSourcePath,
		mdAssetsDestPath: mdAssetsDestPath,
		wantsBannerFirst: wantsBannerForFirst,
	}
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		img, ok := node.(*ast.Image)
		if !ok {
			return ast.GoToNext, false
		}
		if !entering {
			// Everything is written on entering; nothing to do on exit.
			return ast.SkipChildren, true
		}
		return st.renderImage(w, img)
	}
}

func (st *imageHookState) renderImage(w io.Writer, img *ast.Image) (ast.WalkStatus, bool) {
	dest := string(img.Destination)
	alt := plainText(img)
	title := string(img.Title)

	isFirst := !st.seenFirst
	st.seenFirst = true

	cleanPath, directive := ParseDestination(dest)

	absPath, ok := ResolveImageSource(st.mdSourcePath, cleanPath)
	if !ok {
		if IsLocalImagePath(cleanPath) {
			if l := logs.GetLoggerctx(st.ctx); l != nil {
				l.Sugar().Warnf("could not find local image %q referenced from %s - rendering as plain <img> (check the file exists at that path)", cleanPath, st.mdSourcePath)
			}
		}
		writeDefaultImg(w, dest, alt, title)
		return ast.SkipChildren, true
	}

	preset := PresetContent
	switch {
	case directive.Preset != "":
		preset = directive.Preset
	case isFirst && st.wantsBannerFirst:
		preset = PresetBanner
	}
	if preset != PresetContent && !PresetExists(preset) {
		if l := logs.GetLoggerctx(st.ctx); l != nil {
			l.Sugar().Warnf("image preset %q not found in config, falling back to %q for %s", preset, PresetContent, absPath)
		}
		preset = PresetContent
	}

	targetWidth := PresetWidth(preset)
	targetHeight := 0
	if directive.Preset == "" {
		if directive.Width > 0 {
			targetWidth = directive.Width
		} else if directive.Height > 0 {
			targetWidth, targetHeight = 0, directive.Height
		}
	}

	result, err := ProcessImage(st.ctx, absPath, targetWidth, targetHeight, PresetBreakpoints(preset), st.mdAssetsDestPath)
	if err != nil {
		if l := logs.GetLoggerctx(st.ctx); l != nil {
			l.Sugar().Warnf("image processing failed for %s, falling back to plain <img>: %v", absPath, err)
		}
		writeDefaultImg(w, dest, alt, title)
		return ast.SkipChildren, true
	}

	fetchPriority := PresetFetchPriority(preset)
	if directive.FetchPriority != "" {
		fetchPriority = directive.FetchPriority
	}

	writePicture(w, result, alt, title, PresetLoading(preset), fetchPriority, PresetSizes(preset))
	return ast.SkipChildren, true
}

func plainText(node ast.Node) string {
	var sb strings.Builder
	ast.WalkFunc(node, func(n ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if leaf := n.AsLeaf(); leaf != nil {
				sb.Write(leaf.Literal)
			}
		}
		return ast.GoToNext
	})
	return sb.String()
}

// writeDefaultImg reproduces byte-for-byte what gomarkdown's own default
// Image renderer would have written (see gomarkdown/html/renderer.go
// imageEnter/imageExit), using gomarkdown's own exported escaping helpers so
// the fallback is indistinguishable from the pre-feature output.
func writeDefaultImg(w io.Writer, dest, alt, title string) {
	io.WriteString(w, `<img src="`)
	gmhtml.EscLink(w, []byte(dest))
	io.WriteString(w, `" alt="`)
	gmhtml.EscapeHTML(w, []byte(alt))
	if title != "" {
		io.WriteString(w, `" title="`)
		gmhtml.EscapeHTML(w, []byte(title))
	}
	io.WriteString(w, `" />`)
}

func writePicture(w io.Writer, r *model.Result, alt, title, loading, fetchPriority, sizes string) {
	io.WriteString(w, "<picture>")
	if len(r.WebP) > 0 {
		io.WriteString(w, `<source type="image/webp" srcset="`)
		gmhtml.EscapeHTML(w, []byte(srcset(r.WebP)))
		io.WriteString(w, `"`)
		writeSizesAttr(w, sizes)
		io.WriteString(w, ">")
	}

	fallback := r.Fallback
	if len(fallback) == 0 {
		fallback = r.WebP // extremely unlikely, but never leave <img> without a src
	}
	fallbackSrc := ""
	if len(fallback) > 0 {
		fallbackSrc = fallback[len(fallback)-1].URL
	}

	io.WriteString(w, `<img src="`)
	gmhtml.EscLink(w, []byte(fallbackSrc))
	io.WriteString(w, `" width="`)
	fmt.Fprintf(w, "%d", r.Width)
	io.WriteString(w, `" height="`)
	fmt.Fprintf(w, "%d", r.Height)
	io.WriteString(w, `" alt="`)
	gmhtml.EscapeHTML(w, []byte(alt))
	io.WriteString(w, `"`)
	if title != "" {
		io.WriteString(w, ` title="`)
		gmhtml.EscapeHTML(w, []byte(title))
		io.WriteString(w, `"`)
	}
	if loading != "" {
		fmt.Fprintf(w, ` loading="%s"`, loading)
	}
	if fetchPriority != "" && fetchPriority != "auto" {
		fmt.Fprintf(w, ` fetchpriority="%s"`, fetchPriority)
	}
	io.WriteString(w, ` decoding="async"`)
	if len(fallback) > 1 {
		io.WriteString(w, ` srcset="`)
		gmhtml.EscapeHTML(w, []byte(srcset(fallback)))
		io.WriteString(w, `"`)
		writeSizesAttr(w, sizes)
	}
	io.WriteString(w, " />")
	io.WriteString(w, "</picture>")
}

func writeSizesAttr(w io.Writer, sizes string) {
	if sizes != "" {
		fmt.Fprintf(w, ` sizes="%s"`, sizes)
	}
}

func srcset(variants []model.Variant) string {
	parts := make([]string, len(variants))
	for i, v := range variants {
		parts[i] = fmt.Sprintf("%s %dw", v.URL, v.Width)
	}
	return strings.Join(parts, ", ")
}
