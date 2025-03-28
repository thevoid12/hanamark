package testing

import (
	"context"
	"fmt"
	"hanamark/parser"
	"testing"
)

func TestParseMarkdownToHtml(t *testing.T) {
	mdDir := "./test.md"
	//	destDir := "./test.html"
	htlmString, err := parser.ParseMarkdownToHtml(mdDir)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(htlmString)
}

func TestExtractHeadingInMarkdown(t *testing.T) {
	mdDir := "./test.md"
	ctx := context.Background()
	_, err := parser.ExtractHeadingInMarkdown(ctx, mdDir)
	if err != nil {
		t.Errorf(err.Error())
	}
}
