package Crawler

import (
	ds "simpleCrawler/DataStructures"
	"simpleCrawler/Fetcher"
	"sync"
)

// Crawl is a goroutine that will
// Then, it'll send all new urls (not cached) to the outQueue.
func crawl(baseUrl *string, queue *ds.Queue, cache *ds.ConcurrentCache, wg *sync.WaitGroup) {
	go func() {
		for  {
			select {
				case v := <-queue.Out:
					parentUrl := v.(string)
					if _, ok := cache.Load(parentUrl); ok {
						continue
					}
					cache.Store(parentUrl)
					childrenUrls, _ := Fetcher.FetchRelativeUrlsFromPage(baseUrl, parentUrl)
					for _, u := range childrenUrls {
						if _, ok := cache.Load(u); !ok {
							queue.In <- u
						}
					}
			}
		}
	}()
	return
}

func Process(seedUrl *string, maxWorkers int) {
	var cache = ds.NewConcurrentCache()
	queue := ds.NewQueue(maxWorkers)
	wg := &sync.WaitGroup{}
	wg.Add(maxWorkers)
	for i := 0 ; i < maxWorkers; i ++ {
		go crawl(seedUrl, queue, cache, wg)
	}
	queue.In <- *seedUrl
	wg.Wait()
}
