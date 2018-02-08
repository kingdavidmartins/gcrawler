package main

import (
	"fmt"
	"time"
	"github.com/kingdavidmartins/gcrawler/internal/workqueue"
	"sync"
)

func main()  {

	store := make(map[string]string)
	store_mutex := &sync.RWMutex{}
	queue := make(chan string)

	workqueue.SpinUpWorkers(10, queue, store, store_mutex)

	queue <- "https://bost.ocks.org/mike/algorithms/"

	timer := time.NewTimer(time.Second * 5)

	<-timer.C

	fmt.Printf("Enough\n")

	// printf out values in store
	store_mutex.Lock()
	fmt.Println(store)
	store_mutex.Unlock()

}