package workqueue

import (
	"errors"
	"io/ioutil"
	"net/http"
	"log"
	"container/list"
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

func process(url string, buff *list.List) {
	log.Printf("processing %s", url)
	text, err := fetch(url)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse url %s\n", url)
	}

	outgoing_urls, err := parse(text)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse body from %s\n", url)
	}

	for _, u := range outgoing_urls {
		buff.PushBack(u)
	}
}

func pop(lst *list.List) string {
	head := lst.Front()

	lst.Remove(head)
	return head.Value.(string)
}

func worker(id int, queue chan string, store map[string]string) {
	internal_buffer := list.New()

	for url := range queue {
		if store[url] != "" {
			continue
		}

		log.Printf("worker %d received %s from queue", id, url)
		process(url, internal_buffer)

		for {
			if internal_buffer.Len() == 0 {
				break
			}

			next_url := pop(internal_buffer)
			log.Printf("worker %d grabbed %s from internal buffer", id, next_url)
			select {
			case queue <- next_url: 
			  // do nothing
			default: 
				process(next_url, internal_buffer)
			}
		}
	}
}

func SpinUpWorkers(num int, queue chan string, store map[string]string) {
	for id := 0; id < num; id++ {
		go worker(id, queue, store)
	}
}