package crawler

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

var logMsg = log.Infof
var logErr = log.Errorf
var logWarn = log.Warnf
var logFatal = log.Fatalf

// crawls until approximateMaxNodes nodes is reached
func Run(
	endpoint string,
	approximateMaxNodes int32,
	isValidCrawlLink IsValidCrawlLinkFunction,
	connectToDB ConnectToDBFunction,
	addEdgesIfDoNotExist AddEdgeFunction,
	getNewNode GetNewNodeFunction,
) {
	// first connect to db
	if err := connectToDB(); err != nil {
		logFatal("Could not connect do db: %v", err)
	}
	// get starting link if there isn't one already
	if endpoint == "" {
		e, err := getNewNode()
		if err != nil {
			logFatal("Could not find new starting node: %v", err)
		} else {
			endpoint = e
		}
	}
	Crawl(
		endpoint,
		approximateMaxNodes,
		isValidCrawlLink,
		addEdgesIfDoNotExist,
	)
}

// crawls a domain and saves relatives links to a db
func Crawl(
	endpoint string,
	approximateMaxNodes int32,
	isValidCrawlLink IsValidCrawlLinkFunction,
	addEdgesIfDoNotExist AddEdgeFunction,
) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(true),
		colly.CacheDir("/tmp/crawlercache"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})

	// On every a element which has href attribute call callback
	c.OnHTML("html", func(e *colly.HTMLElement) {
		logMsg("parsing %s", e.Request.URL.String())
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
		} else {
			// update metrics
			UpdateMetrics(len(nodesAdded), e.Request.Depth)
		}
		// recurse on new nodes if no stopping condition yet
		if approximateMaxNodes == -1 || totalNodesAdded.get() < approximateMaxNodes {
			for _, url := range nodesAdded {
				err = e.Request.Visit(url)
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
