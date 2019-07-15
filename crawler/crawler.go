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
		link := e.Attr("href")

		if isValidCrawlLink(link) {
			edgeAlreadyExists, err := addEdgeIfDoesNotExist(e.Request.URL.String(), link)
			// only visit link if edge doesnt exist
			if err != nil {
				logMsg("ERROR: %s", err.Error())
			} else if edgeAlreadyExists == false {
				logMsg("added edge %s => %s", e.Request.URL.String(), link)
				e.Request.Visit(link)
			}
		}
	})

	// Start scraping on endpoint
	logMsg("starting at %s", endpoint)
	c.Visit(endpoint)
	// Wait until threads are finished
	c.Wait()
}
