package Crawler

import (
	"fmt"
	"runtime"
	"simpleCrawler/DataStructures"
	"simpleCrawler/Fetcher"
	"sync"
)

func CrawlN(baseUrl *string, urlToParse string, cache *DataStructures.ConcurrentCache, wg *sync.WaitGroup, depth int) {
	go func() {
		if _, ok := cache.Load(urlToParse); ok {
			return
		}
		cache.Store(urlToParse)
		childrenUrls, _ := Fetcher.FetchRelativeUrlsFromPage(baseUrl, urlToParse)
		fmt.Printf("Depth: %v cache size: %v URL %v\n", depth, len(cache.Cache), runtime.NumGoroutine())
		for _, u := range childrenUrls {
			if _, ok := cache.Load(u); !ok {
				wg.Add(1)
				go CrawlN(baseUrl, u, cache, wg, depth + 1)
			}
		}
		wg.Done()
	}()
	return
}

func processN(baseUrl string) {
	var cache = DataStructures.NewConcurrentCache()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go CrawlN(&baseUrl, baseUrl, cache, wg, 0)
	wg.Wait()
}
