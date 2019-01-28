package main

import (
	"fmt"
	"sync"
)

// URLCache is safe to use concurrently.
type URLCache struct {
	v   map[string]fakeResult
	mux sync.Mutex
}

// Adds a URL to URLCache.
func (c *URLCache) Add(url string, result fakeResult) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[url] = result
	c.mux.Unlock()
}

// Value returns the current value for the given URL key.
func (c *URLCache) Value(url string) fakeResult {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()
	return c.v[url]
}

// Exists returns if the URL key exists.
func (c *URLCache) Exists(url string) bool {
	c.mux.Lock()
	exists := false
	// Lock so only one goroutine at a time can access the map c.v.
	if _, ok := c.v[url]; ok {
		exists = true
		/*fmt.Printf("found cache key: %s\n", url)*/
	}
	defer c.mux.Unlock()
	return exists
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, cache URLCache, ch chan string) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	fmt.Printf("In crawl depth: %d url: %s\n", depth, url)
	if depth <= 0 {
		close(ch)
		return
	}
	inCache := cache.Exists(url)
	if inCache {
		/*fmt.Printf("\tthe url is in in cache, skipping: %s\n", url)*/
	} else {
		/*fmt.Printf("\tthe url is not in cache %s\n", url)*/
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			/*fmt.Printf("Error retreiving response from %s : %s\n", url, err)*/
		} else {
			/*fmt.Printf("\tfetched: %s\n", url)*/
			cache.Add(url, fakeResult{body, urls})
			/*fmt.Printf("\tAdded response from %s to cache, there are %d child urls\n", url, len(urls))*/
			for _, u := range urls {
				/*fmt.Printf("\t\tcrawl sub url: %s\n", u)*/
				go Crawl(u, depth-1, fetcher, cache, ch)
			}
		}
	}
	ch <- url
	if depth == 1 {
		fmt.Println("closing channel")
		close(ch)
	}
	return
}

func main() {
	c := URLCache{v: make(map[string]fakeResult)}
	ch := make(chan string)
	go Crawl("https://golang.org/", 4, fetcher, c, ch)
	for s := range ch {
		fmt.Println(s)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
