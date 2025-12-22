package parser

import (
	"context"
	"errors"
	"hanamark/model"
	"hanamark/util"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

// "rss":{
//     "isRssEnabled":true,
//     "title":"void blog",
//     "link":"https://thisisvoid.in",
//     "authorName":"void",
//     "authorEmailID":"thevoidddd1@gmail.com"

//	},
//
// RSS base feed. which is added as part of config if enabled
func GetRssFeed() (feed *feeds.Feed, err error) {
	now := time.Now()
	isRssEnabled := viper.GetBool("rss.isRssEnabled")
	if !isRssEnabled {
		return nil, errors.New("rss not enabled in config")
	}
	link := viper.GetString("rss.link")
	if link == "" {
		return nil, errors.New("link of the root page of your blog is mandatory to config to setup rss")
	}
	//TODO: get these feed details from config
	feed = &feeds.Feed{
		Title: viper.GetString("rss.title"),
		Link:  &feeds.Link{Href: link},
		// Description: "discussion about tech, footie, photos",
		Author:  &feeds.Author{Name: viper.GetString("rss.authorName"), Email: viper.GetString("rss.authorEmailID")},
		Created: now,
	}

	return feed, nil
}

// gets the list of files as rss formatted from the file meta
func GetRssFeedItems(metaList []*model.PageMeta) ([]*feeds.Item, error) {
	items := []*feeds.Item{}
	link := viper.GetString("rss.link")
	if link == "" {
		return nil, errors.New("link of the root page of your blog is mandatory to config to setup rss")
	}

	for _, meta := range metaList {
		link, err := util.JoinURL(link, meta.DestPageDir)
		if err != nil {
			return nil, err
		}
		items = append(items, &feeds.Item{
			Title: meta.PageTitle,
			Link: &feeds.Link{
				Href: link,
				// Rel:    "",
				// Type:   "",
				// Length: "",
			},
			Author:      &feeds.Author{Name: viper.GetString("rss.authorName"), Email: viper.GetString("rss.authorEmailID")},
			Id:          uuid.NewString(),
			Updated:     meta.UpdatedDate,
			Created:     meta.CreatedDate,
			Content:     meta.GenHtml,
			Description: meta.PageTitle, // TODO: we got to figure out a way to extract some subsection from the genetrated content
		})
	}

	return items, nil
}

func GenerateRss(ctx context.Context, feed *feeds.Feed) error {
	rss, err := feed.ToRss()
	if err != nil {
		return err
	}
	// TODO: get the rss path from the config and save the file
	destRootPath := viper.GetString("filepath.destHtmlRoot")
	if destRootPath == "" {
		return errors.New("dest root path in config is empty")
	}
	opName := viper.GetString("rss.rssOutputName")
	if opName == "" {
		opName = "feed.xml"
	}
	if filepath.Ext(opName) != ".xml" {
		opName += ".xml"
	}
	return util.WriteIntoFile(ctx, rss, filepath.Join(destRootPath, opName))
}
