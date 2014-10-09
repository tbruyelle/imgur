package main

import (
	"flag"
	"fmt"
	"github.com/tbruyelle/imgur"
	"io"
	"net/http"
)

var (
	clientId = flag.String("clientId", "", "The Imgur app client id")
	client   *imgur.Client
)

func main() {
	flag.Parse()
	client = imgur.NewClient(*clientId)

	http.HandleFunc("/", randomHandler)
	panic(http.ListenAndServe(":8282", nil))
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	if q == "" {
		q = "cat"
	}
	opt := imgur.SearchOptions{
		Any:  q,
		Type: "gif",
	}
	img, err := client.Random(opt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with the imgur library\n%q", err),
			http.StatusInternalServerError)
		return
	}
	resp, err := http.Get(img.Link)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during download img\n%q", err),
			http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during copy img\n%q", err),
			http.StatusInternalServerError)
		return
	}
}
