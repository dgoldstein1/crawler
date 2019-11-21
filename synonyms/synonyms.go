package synonyms

import (
	"errors"
	"github.com/dgoldstein1/crawler/db"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

// globals
var logErr = log.Errorf
var prefix = "/synonyms/"
var baseEndpoint = "synonyms.com"
var timeout = time.Duration(5 * time.Second)
var c = colly.NewCollector()

// determines if is good link to crawl on
func IsValidCrawlLink(link string) bool {
	validPrefix := strings.HasPrefix(link, "/synonyms/")
	isNotMainPage := strings.ToLower(link) != "/synonyms/main_page"
	noillegalChars := !strings.Contains(link, ":") && !strings.Contains(link, "#")
	return validPrefix && isNotMainPage && noillegalChars
}

// gets random article from metawiki API
// returns article in the form "/synonym/XXXXX"
func GetRandomArticle() (string, error) {
	return "", errors.New("not implemented")
}

// decodes and standaridizes URL
func CleanUrl(link string) string {
	// trim current node if needed
	link = strings.TrimPrefix(link, baseEndpoint)
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
