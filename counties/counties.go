package counties

import (
	"bufio"
	"github.com/dgoldstein1/crawler/db"
	"github.com/dgoldstein1/crawler/util"
	"github.com/dgoldstein1/crawler/wikipedia"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

// globals
var logErr = log.Errorf
var prefix = "/wiki/"
var baseEndpoint = "https://en.wikipedia.org"
var timeout = time.Duration(5 * time.Second)
var counties = make(map[string]bool)

func IsValidCrawlLink(link string) bool {
	// countains the word 'county' in format 'NAME_county,_STATE'
	if !strings.Contains(strings.ToLower(link), "_county,_") {
		return false
	}
	// ensure that is on master list
	return stringInFile(strings.TrimPrefix(link, prefix))
}

func stringInFile(link string) bool {
	// read in counties list if does not exist
	if len(counties) == 0 {
		file, err := os.Open(os.Getenv("COUNTIES_LIST"))
		defer file.Close()
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			counties[strings.ToLower(scanner.Text())] = true
		}
	}
	// now check if string exists in file
	return counties[strings.ToLower(link)]
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
	// first find h3 with id^="Adjacent_couties"
	e.DOM = e.DOM.
		Find(`[id^='Adjacent_counties']`).
		Parent().
		NextUntil("h3").
		Not("p")

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
