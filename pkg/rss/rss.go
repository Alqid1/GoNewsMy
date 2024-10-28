package rss

import (
	"GoNewsMy/pkg/db"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Chanel  RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

func FetchRSSFeed(url string) ([]db.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f RSSFeed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	var data []db.Post
	for _, item := range f.Chanel.Items {
		var p db.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.Link = item.Link

		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		data = append(data, p)
	}
	return data, nil
}
