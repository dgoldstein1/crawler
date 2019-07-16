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
	addEdgesIfDoNotExist AddEdgeFunction,
) {
	err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// Instantiate default collector
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(true),
		colly.CacheDir("/tmp/crawlercache"),
	)
	c.Limit(&colly.LimitRule{Parallelism: 10})

	// On every a element which has href attribute call callback
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// loop through all href attributes adding links
		validURLs := []string{}
		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			// add links which match the schema
			link := e.Attr("href")
			if isValidCrawlLink(link) {
				validURLs = append(validURLs, link)
			}
		})
		// add new nodes to current request URL
		nodesAdded, err := addEdgesIfDoNotExist(e.Request.URL.String(), validURLs)
		if err != nil {
			logMsg("ERROR: %s", err.Error())
		} else {
			// recurse on new nodes
			for _, url := range nodesAdded {
				c.Visit("https://en.wikipedia.org" + url)
			}
		}
	})

	// Start scraping on endpoint
	logMsg("starting at %s", endpoint)
	c.Visit(endpoint)
	// Wait until threads are finished
	c.Wait()
}
