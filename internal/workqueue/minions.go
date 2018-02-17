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

func worker(id int, queue chan string, store map[string]string, store_mutex *sync.RWMutex) {
	internal_buffer := list.New()

	for url := range queue {
		store_mutex.Lock()
		if store[url] != "" {
			store_mutex.Unlock()
			continue
		}
		store_mutex.Unlock()

		log.Printf("worker %d received %s from queue", id, url)
		process(url, internal_buffer, store, store_mutex)

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
				process(next_url, internal_buffer, store, store_mutex)
			}
		}
	}
}

func SpinUpWorkers(num int, queue chan string, store map[string]string, store_mutex *sync.RWMutex) {
	for id := 0; id < num; id++ {
		go worker(id, queue, store, store_mutex)
	}
}
