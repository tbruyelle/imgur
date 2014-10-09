package main

import (
	"flag"
	"fmt"
	"github.com/tbruyelle/imgur"
)

var (
	clientId = flag.String("clientId", "", "The Imgur client id")
)

func main() {
	flag.Parse()
	c := imgur.NewClient(*clientId)

	s, resp, err := c.Search(imgur.SearchOptions{All: "cat", Type: "gif"})
	if err != nil {
		fmt.Println(resp, err)
	} else {
		fmt.Printf("%+v\n", s)
	}
}
