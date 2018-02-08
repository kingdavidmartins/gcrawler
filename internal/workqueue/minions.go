package workqueue

import (
	"errors"
	"io/ioutil"
	"net/http"
	"log"
)

func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("SH!T GOT REAL ~ Couldn't get url")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("SH!T GOT REAL ~ Couldn't read body url")
	}

	return string(body), nil
}

func parse(text string) ([]string, error) {
	outgoing_urls := make([]string, 0)
	return outgoing_urls, nil
}

func process(url string, buff []string) {
	text, err := fetch(url)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse url %s\n", url)
	}

	outgoing_urls, err := parse(text)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse body from %s\n", url)
	}

	for _, u := range outgoing_urls {
		buff = append(buff, u)
	}
}