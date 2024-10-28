package rss

import (
	"testing"
)

func TestFetchRSSFeed(t *testing.T) {
	feed, err := FetchRSSFeed("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(feed), feed)
}
