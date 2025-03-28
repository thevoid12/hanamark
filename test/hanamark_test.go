package testing

import (
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
