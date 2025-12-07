package model

import "time"

// data for each page. this is a flat structure
// so can be customized based on usecase
type PageMeta struct {
	GenHtml     string // generated html
	ReadTime    int    // total time to read that page
	PageName    string
	PageTitle   string
	Date        time.Time
	DestPageDir string
	BaseFile    string // base file is the index file for each subfiles if exists or the root file itself eg blogs.html,index.html,projects.html etc
	// TODO: remove the basefile variable

}

type ListPage struct {
	Base     string
	TempPath string
}
