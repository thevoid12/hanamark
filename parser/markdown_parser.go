package parser

import (
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// ParseMarkdownToHtml parses markdown from source directory into html string
func ParseMarkdownToHtml(sourceMDPath string) (string, error) {

	mdInputfile, err := os.ReadFile(sourceMDPath)
	if err != nil {
		return "", err
	}
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.SuperSubscript | parser.Includes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdInputfile)
	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.CompletePage
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	result := markdown.Render(doc, renderer)

	return string(result), nil
}
