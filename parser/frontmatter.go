package parser

import (
	"bytes"
	"context"
	logs "hanamark/logger"
	"os"

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

	f, err := os.Open(FilePath)
	if err != nil {
		return nil, err
	}
	fm = map[string]any{} // dynamic key/value map

	// Parse frontmatter into the map; rest contains Markdown without front matter
	_, err = frontmatter.Parse(f, &fm)
	if err != nil {
		return nil, err
	}
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
