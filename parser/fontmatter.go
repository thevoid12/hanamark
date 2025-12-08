package parser

import (
	"context"
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
