package crawler

import (
	"github.com/gocolly/colly"
	"log"
)

var logMsg = log.Printf
// crawls a domain and saves relatives links to a db
func Crawl(
	endpoint string,
	maxDepth int,
	isValidCrawlLink IsValidCrawlLinkFunction,
	connectToDB ConnectToDBFunction,
	addEdgeIfDoesNotExist AddEdgeFunction,
) {
	err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// Instantiate default collector
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{Parallelism: 10})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

	})

	// Start scraping on endpoint
	logMsg("starting at %s", endpoint)
	c.Visit(endpoint)
	// Wait until threads are finished
	c.Wait()
}
