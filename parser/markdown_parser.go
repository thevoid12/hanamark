package parser

import (
	"context"
	"fmt"
	"hanamark/util"
	"os"

	"github.com/gomarkdown/markdown/ast"

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
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	result := markdown.Render(doc, renderer)

	return string(result), nil
}

func ExtractHeadingInMarkdown(ctx context.Context, sourceMDPath string) (string, error) {
	// Read the Markdown file
	mdInputfile, err := os.ReadFile(sourceMDPath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Configure parser extensions
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoEmptyLineBeforeBlock |
		parser.SuperSubscript |
		parser.Includes

	// Create parser and parse document
	p := parser.NewWithExtensions(extensions)
	parentDoc := p.Parse(mdInputfile)
	firstChild := parentDoc.AsContainer().Children[0]
	heading := RecurseThroughAST(firstChild)
	heading = util.CleanSpaces(heading)

	return heading, nil
}

// RecurseThroughAST recurses through the abstract syntax tree to fetch the heading
func RecurseThroughAST(node ast.Node) string {
	if node == nil {
		return ""
	}

	// If it's a leaf node, extract and return its text
	leafNode := node.AsLeaf()
	if leafNode != nil {
		return string(leafNode.Literal)
	}

	// traversing all child nodes and accumulating text
	var result string
	for _, child := range node.GetChildren() {
		result += RecurseThroughAST(child) + " "
	}

	return result
}
