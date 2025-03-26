package testing

import (
	"hanamark/parser"
	"testing"
)

func TestParseMarkdownToHtml(t *testing.T) {
	mdstring := "# this is a test heading - void "
	err := parser.ParseMarkdownToHtml(mdstring)
	if err != nil {
		t.Errorf(err.Error())
	}
}
