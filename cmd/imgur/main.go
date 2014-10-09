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

	s, resp, err := c.Search("cat ext:gif")
	if err != nil {
		fmt.Println(resp, err)
	} else {
		fmt.Println(s)
	}
}
