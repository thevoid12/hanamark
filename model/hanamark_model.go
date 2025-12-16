package model

import "time"

// data for each page. this is a flat structure
// so can be customized based on usecase
type PageMeta struct {
	GenHtml        string // generated html
	ReadTime       int    // total time to read that page
	PageName       string
	PageTitle      string
	Date           time.Time
	DestPageDir    string
	FrontMatterMap map[string]any // will have all the frontmatter from the md file
	BaseFile       string         // base file is the index file for each subfiles if exists or the root file itself eg blogs.html,index.html,projects.html etc
	// TODO: remove the basefile variable
	Tags []*Tag
}

type ListPage struct {
	Base     string
	TempPath string
}

// Constants for All the front matter which are supported by the system
const (
	TEMPLATE = "template" // custom templating instead of single and list if needed
	TAGS     = "tags"     // TODO: needs to be implemented
	DRAFT    = "draft"    // draft mode documents arent
)

// Tags
type Tag struct {
	TagName         string // this is the path where the list of tag html is generated ie /pointB/tags/tagname.html
	TagDestPath     string
	TagTemplatePath string // the index.html template path for each unique tag. If it doesnt exists we take a generic index.html in tags
	FileHeading     string // heading of the file which is using the tag
	FileDestPath    string // the path of the md->html comverted file which has the tag
}

var TagMap map[string]Tag // key is the tag name value is the tag property

type TagList struct {
	TagName     string
	TagDestPath string // the path that leads us to the individual tag list
	Count       int    // the number of records with the tag name
}
