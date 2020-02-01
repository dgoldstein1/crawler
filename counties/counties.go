package counties

import (
	"github.com/dgoldstein1/crawler/db"
	"github.com/dgoldstein1/crawler/util"
	"github.com/dgoldstein1/crawler/wikipedia"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// globals
var logErr = log.Errorf
var prefix = "/wiki/"
var baseEndpoint = "https://en.wikipedia.org"
var timeout = time.Duration(5 * time.Second)

func IsValidCrawlLink(link string) bool {
	// countains the word 'county' in format 'NAME_county,_STATE'
	if !strings.Contains(strings.ToLower(link), "_county,_") {
		return false
	}
	// assert that only contains one ',_' (more than one denotes town)
	if strings.Count(link, ",_") > 1 {
		return false
	}
	// national registry
	if strings.Contains(strings.ToLower(link), "national_register_of_historic_places") {
		return false
	}
	return wikipedia.IsValidCrawlLink(link)
}

func GetRandomNode() (string, error) {
	return util.ReadRandomLineFromFile(
		"COUNTIES_LIST",
		baseEndpoint,
		prefix,
		false,
	)
}

// decodes and standaridizes URL
func CleanUrl(link string) string {
	return wikipedia.CleanUrl(link)
}

// filters down full page body to elements we want to focus on
func FilterPage(e *colly.HTMLElement) (*colly.HTMLElement, error) {
	e.DOM = e.DOM.Find("#bodyContent")
	return e, nil
}

// adds edge to DB, returns new neighbors added (to crawl on)
func AddEdgesIfDoNotExist(
	currentNode string,
	neighborNodes []string,
) (
	neighborsAdded []string,
	err error,
) {
	return db.AddEdgesIfDoNotExist(
		currentNode,
		neighborNodes,
		CleanUrl,
		baseEndpoint,
	)
}
