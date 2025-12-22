package parser

import (
	"bytes"
	"context"
	"fmt"
	logs "hanamark/logger"
	"hanamark/model"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
)

// font matter is the text/meta data you write on top
// for example:
//
// ---
// title: "test"
// date: "2022-08-23"
// author: "vinod"
// description: "this is a test"
// tags: ["test", "test2"]
// ---
// able to add font matter makes the parser more powerful and dynamic
func ParseFrontMatter(ctx context.Context, FilePath string) (fm map[string]any, err error) {
	if FilePath == "" {
		return nil, fmt.Errorf("failed to parse frontmatter: file path is empty")
	}

	f, err := os.Open(FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s for frontmatter parsing: %w", FilePath, err)
	}
	fm = map[string]any{} // dynamic key/value map

	// Parse frontmatter into the map; rest contains Markdown without front matter
	_, err = frontmatter.Parse(f, &fm)
	if err != nil {
		return nil, err
	}

	lower := make(map[string]any, len(fm)) // pre-size: small win
	for k, v := range fm {
		lower[strings.ToLower(k)] = v
	}
	fm = lower
	return fm, nil
}

// we have to remove the frontmatter from the md else it will also be displayed along the html
func StripFrontMatter(ctx context.Context, data []byte) []byte {
	rest := data
	l := logs.GetLoggerctx(ctx)

	// skip leading blank/whitespace-only lines
	for {
		lineEnd := bytes.Index(rest, []byte("\n"))
		if lineEnd == -1 {
			return data
		}
		line := rest[:lineEnd]
		if len(bytes.TrimSpace(line)) == 0 {
			rest = rest[lineEnd+1:]
			continue
		}
		break
	}

	// check if first non-empty line is ---
	if !bytes.Equal(bytes.TrimSpace(rest[:bytes.Index(rest, []byte("\n"))]), []byte("---")) {
		return data
	}

	// remove opening line
	rest = rest[bytes.Index(rest, []byte("\n"))+1:]

	// look for closing ---
	for {
		lineEnd := bytes.Index(rest, []byte("\n"))
		if lineEnd == -1 {
			l.Sugar().Error("marformed front matter. Plese check the source markdown file")
			return data // malformed; do nothing
		}

		line := rest[:lineEnd]
		if bytes.Equal(bytes.TrimSpace(line), []byte("---")) {
			// return content after closing marker
			return rest[lineEnd+1:]
		}

		rest = rest[lineEnd+1:]
	}
}

// check for a paticular font matter type exist in fontmatter
// add the check into this whenever a new fontmatter key is added
// you can either provide the filepath or if you have the frontmatter if you have it already
func FrontMatterValidator(ctx context.Context, FilePath string, fm map[string]any, fmKey model.FrontMatterKey) (isPresent bool, value any, err error) {
	if len(fm) == 0 {
		if FilePath == "" {
			return false, nil, fmt.Errorf("cannot validate frontmatter key %s: both frontmatter map and file path are empty", fmKey)
		}
		fm, err = ParseFrontMatter(ctx, FilePath)
		if err != nil {
			return false, nil, fmt.Errorf("failed to validate frontmatter: %w", err)
		}
	}
	v, ok := fm[strings.ToLower(string(fmKey))]
	if !ok {
		return false, nil, nil
	}
	return true, v, nil
}
