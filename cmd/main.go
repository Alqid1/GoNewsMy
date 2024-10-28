package main

import (
	"GoNewsMy/pkg/api"
	"GoNewsMy/pkg/db"
	"GoNewsMy/pkg/rss"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const constr string = "host=localhost port=5432 user=postgres password=postgres dbname=GoNewsMy sslmode=disable"

type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	datb, err := db.New(constr)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(datb)

	b, err := ioutil.ReadFile("cmd/config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	chRSS := make(chan db.Post)
	chErr := make(chan error)

	for _, url := range config.URLS {
		go parseRSS(url, datb, config.Period, chRSS, chErr)
	}

	go func() {
		for post := range chRSS {
			datb.NewPost(post)
		}
	}()

	go func() {
		for err := range chErr {
			log.Println("ошибка:", err)
		}
	}()

	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

func parseRSS(url string, datb *db.DatBase, period int, postsRSS chan<- db.Post, errors chan<- error) {
	for {
		posts, err := rss.FetchRSSFeed(url)
		if err != nil {
			errors <- err
			continue
		}
		for _, p := range posts {
			postsRSS <- p
		}
		time.Sleep(time.Minute * time.Duration(period))
	}
}
