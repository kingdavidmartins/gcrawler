package workqueue

import (
	"container/list"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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
	outgoingUrls := make([]string, 0)
	return outgoingUrls, nil
}

func process(url string, buff *list.List, store map[string]string, storeMutex *sync.RWMutex) {
	log.Printf("processing %s", url)
	text, err := fetch(url)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse url %s\n", url)
	}

	outgoingUrls, err := parse(text)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse body from %s\n", url)
	}

	// push into store
	storeMutex.Lock()
	store[url] = url
	storeMutex.Unlock()

	for _, u := range outgoingUrls {
		buff.PushBack(u)
	}
}

func pop(lst *list.List) string {
	head := lst.Front()

	lst.Remove(head)
	return head.Value.(string)
}
