package model

import "time"

// data for each page. this is a flat structure
// so can be customized based on usecase
type PageMeta struct {
	GenHtml     string // generated html
	PageName    string
	PageTitle   string
	PageType    string
	Date        time.Time
	DestPageDir string
}

// key is baseName value is  pageMeta
var MiscData map[string][]*PageMeta
