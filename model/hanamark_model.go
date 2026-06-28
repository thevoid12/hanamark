package model

import "time"

// data for each page. this is a flat structure
// so can be customized based on usecase
type PageMeta struct {
	GenHtml        string // generated html
	ReadTime       int    // total time to read that page
	PageName       string
	PageTitle      string
	CreatedDate    time.Time
	UpdatedDate    time.Time
	DestPageDir    string
	FrontMatterMap map[string]any // will have all the frontmatter from the md file
	BaseFile       string         // base file is the index file for each subfiles if exists or the root file itself eg blogs.html,index.html,projects.html etc
	Tags           []*Tag
}

type ListPage struct {
	Base     string
	TempPath string
}

// after parsing the frontmatter it returns a map of this type
type FrontMatter map[string]any

// Constants for All the front matter which are supported by the system
// FrontMatterKey represents a supported front matter field
type FrontMatterKey string

// Supported front matter keys
const (
	TEMPLATE FrontMatterKey = "template"   // custom templating instead of single and list if needed
	TAGS     FrontMatterKey = "tags"       // tag list for the document
	DRAFT    FrontMatterKey = "draft"      // draft mode documents are not published
	DATE       FrontMatterKey = "created_on" // date at which the md file is first created
	UPDATED_ON FrontMatterKey = "updated_on" // date at which the md file was last updated; falls back to created_on when absent
	RSS        FrontMatterKey = "rss"        // rss support
	TITLE      FrontMatterKey = "title"      // custom title for blogs (using this title overwrites the actual generated title in list pages of blogs)
)

// Tags
type Tag struct {
	TagName         string // this is the path where the list of tag html is generated ie /pointB/tags/tagname.html
	TagDestPath     string
	TagTemplatePath string // the index.html template path for each unique tag. If it doesnt exists we take a generic index.html in tags
	FileHeading     string // heading of the file which is using the tag
	FileDestPath    string // the path of the md->html comverted file which has the tag
}

type TagList struct {
	TagName     string
	TagDestPath string // the path that leads us to the individual tag list
	Count       int    // the number of records with the tag name
}

// IndexHomepageHtml type constants for config
const (
	IndexTypeSection = "section"
	IndexTypePage    = "page"
)

var TimeLayouts = []string{
	time.RFC3339,     // 2006-01-02T15:04:05Z07:00
	time.RFC3339Nano, // 2006-01-02T15:04:05.999999999Z07:00

	// Date only
	"2006-01-02", // 2024-06-01
	"2006/01/02", // 2024/06/01

	// Date + time (no timezone)
	"2006-01-02 15:04",    // 2024-06-01 10:30
	"2006-01-02 15:04:05", // 2024-06-01 10:30:45

	// Human-friendly formats
	"02-01-2006",      // 01-06-2024 (DD-MM-YYYY)
	"02/01/2006",      // 01/06/2024
	"02 Jan 2006",     // 01 Jun 2024
	"02 January 2006", // 01 June 2024

	// With time + words
	"02 Jan 2006 15:04",     // 01 Jun 2024 10:30
	"02 January 2006 15:04", // 01 June 2024 10:30

	// Filename-style dates (common in SSGs)
	"20060102", // 20240601
}
