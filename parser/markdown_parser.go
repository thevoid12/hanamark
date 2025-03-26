package parser

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func ParseMarkdownToHtml(mdFilepath string) error {

	mdInputfile, err := os.ReadFile(mdFilepath)
	if err != nil {
		return err
	}
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdInputfile)
	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	result := markdown.Render(doc, renderer)
	f, err := os.Create("../test/test.html")
	if err != nil {
		fmt.Println(err)

	}
	defer f.Close()
	_, err = f.Write(result)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
