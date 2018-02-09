package workqueue

import (
	"errors"
	"net/http"
	"log"
	"container/list"
	"strings"
	"sync"
	"golang.org/x/net/html"
	"fmt"
)

func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("SH!T GOT REAL ~ Couldn't get url")
	}

	// close the response body when finished with it
	defer resp.Body.Close()

	body := resp.Body

	if err != nil {
		return "", errors.New("SH!T GOT REAL ~ Couldn't read body url")
	}

	z := html.NewTokenizer(body)

	for {
			tt := z.Next()
	
			switch {
			case tt == html.ErrorToken:
				// End of the document, we're done
					// return nothing
			case tt == html.StartTagToken:
					t := z.Token()

					for _, a := range t.Attr {
						not_http := strings.Index(a.Val, "http") == 0
						if !not_http {
							continue
						}
						if a.Key == "href" {
							fmt.Println(a.Val)
						}
					}
			}
	}

	// return string(body), nil
}

func parse(text string) ([]string, error) {
	outgoing_urls := make([]string, 0)
	return outgoing_urls, nil
}

func process(url string, buff *list.List, store map[string]string, store_mutex *sync.RWMutex) {
	log.Printf("processing %s", url)
	text, err := fetch(url)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse url %s\n", url)
	}

	outgoing_urls, err := parse(text)
	if err != nil {
		log.Fatalf("SH!T GOT REAL ~ Failed to parse body from %s\n", url)
	}

	// push into store
	store_mutex.Lock()
	store[url] = url
	store_mutex.Unlock()

	for _, u := range outgoing_urls {
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