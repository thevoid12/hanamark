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
