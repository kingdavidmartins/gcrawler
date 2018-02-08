package workqueue

import (
	"errors"
	"io/ioutil"
	"net/http"
	// "log"
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