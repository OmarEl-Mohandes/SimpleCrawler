package Crawler

import (
	"log"
	"os"
	ds "simpleCrawler/DataStructures"
	"simpleCrawler/Fetcher"
	"sync"
	"time"
)

const (
	maxWorkerIdleDuration = time.Second * 5
)

// Crawl is a goroutine that will act as worker that consumes from the queue given (using In Channel).
// It'll mark the url in the cache if not seen.
// Then, it fetches the relative urls based on the baseUrl given, and it'll send all new urls (not cached)
// to the Queue again (using Out channel).
// If the worker is Idle for maxWorkerIdleDuration, then the worker will quit.
func crawl(baseUrl *string, queue *ds.Queue, cache *ds.ConcurrentCache, wg *sync.WaitGroup, logger *log.Logger) {
	go func() {
		idleTimer := time.NewTimer(maxWorkerIdleDuration)
		resetTimer := func() {
			if !idleTimer.Stop() {
				<-idleTimer.C
			}
			idleTimer.Reset(maxWorkerIdleDuration)
		}
		for  {
			select {
				case v := <-queue.Out:
					resetTimer()
					parentUrl := v.(string)
					if _, ok := cache.Load(parentUrl); ok {
						continue
					}
					logger.Printf("URL found: %v\n", parentUrl)
					cache.Store(parentUrl)
					childrenUrls, _ := Fetcher.FetchRelativeUrlsFromPage(baseUrl, parentUrl)
					for _, u := range childrenUrls {
						if _, ok := cache.Load(u); !ok {
							queue.In <- u
						}
					}
				case <-idleTimer.C:
					wg.Done()
					return
			}
		}
	}()
	return
}

func Process(seedUrl *string, maxWorkers int) {
	logger := log.New(os.Stdout, "", 0)
	var cache = ds.NewConcurrentCache()
	queue := ds.NewQueue(maxWorkers)
	wg := &sync.WaitGroup{}
	wg.Add(maxWorkers)
	for i := 0 ; i < maxWorkers; i ++ {
		go crawl(seedUrl, queue, cache, wg, logger)
	}
	queue.In <- *seedUrl
	wg.Wait()
}
