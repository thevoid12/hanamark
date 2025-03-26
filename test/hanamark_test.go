package testing

import (
	"hanamark/parser"
	"testing"
)

func TestParseMarkdownToHtml(t *testing.T) {
	mdDir := "./test.md"
	err := parser.ParseMarkdownToHtml(mdDir)
	if err != nil {
		t.Errorf(err.Error())
	}
}
