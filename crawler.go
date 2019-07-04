package crawler

import (
	"github.com/gocolly/colly"
	"regexp"
	"log"
	"os"
	"net/http"
	"io/ioutil"
)

// crawls a domain and saves relatives links to a db
func Crawl(endpoint string, urlRegex *regexp.Regexp, maxDepth int) {
	err := connectToDB()
	if err != nil {
		log.Fatal("Could not connect to DB")
		panic(err)
	}
	// Instantiate default collector
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{Parallelism: 10})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		if urlRegex.MatchString(link) {
			neighborExists, _ := addToDB(e.Request.URL.String(), link)
			if !neighborExists {
				e.Request.Visit(link)
			}
		}
	})

	// Start scraping on endpoint
	c.Visit(endpoint)
	// Wait until threads are finished
	c.Wait()
}

// adds edge to DB, returns (true) if neighbor node exists
func addToDB(currentNode string, neighborNode string) (bool, error) {
	// os.Getenv("GRAPH_DB_ENDPOINT")
	return true, nil
}

// connects to given databse
func connectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT") + "/metrics")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	// handling error and doing stuff with body that needs to be unit tested
	return err
}
