package crawler

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var logMsg = log.Infof
var logErr = log.Errorf
var logWarn = log.Warnf

// crawls a domain and saves relatives links to a db
func Crawl(
	endpoint string,
	approximateMaxNodes int32,
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
		colly.Async(true),
		colly.CacheDir("/tmp/crawlercache"),
	)
	c.Limit(&colly.LimitRule{Parallelism: 10})

	nodesVisited := asyncInt(0)

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
			logErr("error adding '%s': %s", e.Request.URL.String(), err.Error())
			return
		}
		// check stopping condition
		nodesVisited.incr(int32(len(nodesAdded)))
		// recurse on new nodes if no stopping condition yet
		if approximateMaxNodes == -1 || nodesVisited.get() < approximateMaxNodes {
			for _, url := range nodesAdded {
				err = c.Visit(url)
				if err != nil {
					logWarn("Error visiting '%s', %v", url, err)
				}
			}
		}
	})

	// Start scraping on endpoint
	logMsg("starting at %s", endpoint)
	c.Visit(endpoint)
	// Wait until threads are finished
	c.Wait()
}

// increments async int by "n"
func (c *asyncInt) incr(n int32) int32 {
	return atomic.AddInt32((*int32)(c), n)
}

// decrement astnc int
func (c *asyncInt) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
