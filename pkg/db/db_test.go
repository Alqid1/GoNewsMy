package db

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestDatBase_Posts(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := Post{
		Title: "Test Post",
		Link:  strconv.Itoa(rand.Intn(1_000_000_000)),
	}
	var constr string = "host=localhost port=5432 user=postgres password=postgres dbname=GoNewsMy sslmode=disable"
	dat, err := New(constr)
	if err != nil {
		t.Fatal(err)
	}
	err = dat.NewPost(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := dat.Posts(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}
