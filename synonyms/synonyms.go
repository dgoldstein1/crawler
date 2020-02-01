package synonyms

import (
	"github.com/dgoldstein1/crawler/db"
	"github.com/dgoldstein1/crawler/util"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

// globals
var logErr = log.Errorf
var prefix = "/synonym/"
var baseEndpoint = "http://www.synonyms.com"
var timeout = time.Duration(5 * time.Second)
var c = colly.NewCollector()

// determines if is good link to crawl on
func IsValidCrawlLink(link string) bool {
	validPrefix := strings.HasPrefix(link, prefix)
	noillegalChars := !strings.Contains(link, ":") && !strings.Contains(link, "#")
	valid := validPrefix && noillegalChars
	if !valid {
		logErr("invalid link found %s. validPrefix : %v, noillegalChars: %v", link, validPrefix, noillegalChars)
	}
	return valid
}

func GetRandomNode() (string, error) {
	return util.ReadRandomLineFromFile(
		"ENGLISH_WORD_LIST_PATH",
		baseEndpoint,
		prefix,
	)
}

// decodes and standaridizes URL
func CleanUrl(link string) string {
	// trim current node if needed
	link = strings.TrimPrefix(link, baseEndpoint)
	// also replace 'https'-- synonyms.com upgrades connection
	httpsPrefix := strings.ReplaceAll(baseEndpoint, "http", "https")
	link = strings.TrimPrefix(link, httpsPrefix)
	link = strings.TrimPrefix(link, prefix)
	link = strings.ToLower(link)
	link = strings.ReplaceAll(link, "_", " ")
	// decode string
	link, err := url.QueryUnescape(link)
	if err != nil {
		logErr("Could not decode string %s: %v", link, err)
	}
	return link
}

// filters down full page body to elements we want to focus on
func FilterPage(e *colly.HTMLElement) (*colly.HTMLElement, error) {
	e.DOM = e.DOM.Find(".syns")
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
