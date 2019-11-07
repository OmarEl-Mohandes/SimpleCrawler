package main

import (
	"flag"
	"fmt"
	"os"
	"simpleCrawler/Crawler"
	"time"
)

func parseInput()(seedUrl string, maxWorkers int, durationInSeconds int) {
	flag.StringVar(&seedUrl, "s", "", "Seed url to start crawling from.")
	flag.IntVar(&maxWorkers, "w", 1000, "Max number of workers to crawl.")
	flag.IntVar(&durationInSeconds, "d", -1, "Number of seconds to crawl, default will be forever until no more crawling is needed")
	flag.Parse()
	if seedUrl == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	return
}

func main() {
	seedUrl, maxWorkers, durationInSeconds := parseInput()
	crawlManager := Crawler.NewCrawlManger(seedUrl, maxWorkers)
	now := time.Now()
	crawlManager.Start()
	if durationInSeconds == -1 {
		crawlManager.WaitUntilWorkersDone()
	} else {
		time.Sleep(time.Second * time.Duration(durationInSeconds))
		crawlManager.Shutdown()
	}
	fmt.Printf("Crawling took %v\n", time.Now().Sub(now))
}

