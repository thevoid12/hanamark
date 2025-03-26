package testing

import (
	"hanamark/parser"
	"testing"
)

func TestParseMarkdownToHtml(t *testing.T) {
	mdDir := "./test.md"
	destDir := "./test.html"
	err := parser.ParseMarkdownToHtml(mdDir, destDir)
	if err != nil {
		t.Errorf(err.Error())
	}
}
