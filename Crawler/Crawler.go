package Crawler

import (
	"fmt"
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

type CrawlManager struct {
	queue 		 *ds.Queue
	cache 		 *ds.ConcurrentCache
	wg    		 *sync.WaitGroup
	logger  	 *log.Logger
	baseUrl 	 *string
	numOfWorkers int
}

func NewCrawlManger(seedUrl string, numOfWorkers int) *CrawlManager {
	crawlManager := &CrawlManager{
		queue:   ds.NewQueue(numOfWorkers),
		cache:   ds.NewConcurrentCache(),
		wg:      &sync.WaitGroup{},
		baseUrl: &seedUrl,
		logger:  log.New(os.Stdout, "", 0),
		numOfWorkers: numOfWorkers,
	}
	for i := 0 ; i < numOfWorkers ; i ++ {
		crawlManager.launchWorker()
	}
	return crawlManager
}

func (cr *CrawlManager) Start()  {
	cr.queue.In<- *cr.baseUrl
}

func (cr *CrawlManager) WaitUntilWorkersDone()  {
	cr.wg.Wait()
}

func (cr *CrawlManager) Shutdown() {
	defer func() {
		// Recovering from the negative sync number.
		recover()
		fmt.Println("Done!")
	}()
	cr.logger.Printf("Sutting Down %v workers..\n", cr.numOfWorkers)
	for i := 0 ; i < cr.numOfWorkers ; i ++ {
		cr.wg.Done()
	}
}

// Worker is a goroutine that consumes from the queue given (using In Channel).
// It'll mark the url in the cache if not seen.
// Then, it fetches the relative urls based on the baseUrl given, and it'll send all new urls (not cached)
// to the Queue again (using Out channel).
// If the worker is Idle for maxWorkerIdleDuration, then the worker will quit.
func (cr *CrawlManager) launchWorker() {
	cr.wg.Add(1)
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
				case v := <-cr.queue.Out:
					resetTimer()
					parentUrl := v.(string)
					if _, ok := cr.cache.Load(parentUrl); ok {
						continue
					}
					cr.logger.Printf("URL found: %v\n", parentUrl)
					cr.cache.Store(parentUrl)
					childrenUrls, _ := Fetcher.FetchRelativeUrlsFromPage(cr.baseUrl, parentUrl)
					for _, u := range childrenUrls {
						if _, ok := cr.cache.Load(u); !ok {
							cr.queue.In <- u
						}
					}
				case <-idleTimer.C:
					cr.wg.Done()
					return
			}
		}
	}()
	return
}
