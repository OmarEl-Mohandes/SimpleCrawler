## Overview

This is a simple crawler program written in Go. Using one url seed, it should fetch all relative urls recursively for the domain, and will print it on the console.

### Requirements for building the program

- You'll need `go1.13` to build.
- To build:  `go build`.
- To test: `go test` . 


#### Usage

```
$ ./simpleCrawler -s "https://example.com" -w 5000 -w 60

Usage of ./simpleCrawler:
  -d int
    	Number of seconds to crawl, default will be forever until no more crawling is needed (default -1)
  -s string
    	Seed url to start crawling from. (required)
  -w int
    	Max number of workers to crawl. (default 1000)
```

### Todo
This project isn't "production ready" due to couple of points:

- Write unit tests for the Crawler & Fetcher packages using [gomock](https://github.com/golang/mock).
- Implement politeness delay to not get throttled or cause pain for people.
- Make the number of cores ```GOMAXPROCS``` configurable, as currently it's ```1```, depending on your number of logical available cores.
- Sometimes this program might hit ```socket: too many open files``` if you use a lot of workers (depends on your default limits).
    - Your limit for max open files per process (```ulimit -n``` to see the current limit).
- Handling errors for gracefully: 
    - e.g Currently, Fetcher will swallow the errors if fetching url failed (will assume this url has no children). We should implement retries with exponential backoff.
- Writing the results to a file, with status update every (e.g 10 seconds) on the console. 
    - Currently, it logs every url found on the console. (you can ``` ./simpleCrawler .. > urls.txt ```)

