package main

import (
	"flag"
	"fmt"
	"os"
	"simpleCrawler/Crawler"
	"time"
)

func parseInput()(seedUrl string, maxWorkers int) {
	flag.StringVar(&seedUrl, "s", "", "Seed url to start crawling from.")
	flag.IntVar(&maxWorkers, "w", 1000, "Max number of workers to crawl.")
	flag.Parse()
	if seedUrl == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	return
}

func main() {
	seedUrl, maxWorkers := parseInput()
	now := time.Now()
	Crawler.Process(&seedUrl, maxWorkers)
	fmt.Println(time.Now().Sub(now))
}

