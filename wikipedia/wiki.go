package wikipedia

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgoldstein1/crawler/db"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// globals
var logErr = log.Errorf
var prefix = "/wiki/"
var baseEndpoint = "https://en.wikipedia.org"
var metawikiEndpoint = "https://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exlimit=max&explaintext&exintro&generator=random&grnnamespace=0&grnlimit=1ts="
var timeout = time.Duration(5 * time.Second)
var c = colly.NewCollector()

// determines if is good link to crawl on
func IsValidCrawlLink(link string) bool {
	validPrefix := strings.HasPrefix(link, "/wiki/")
	isNotMainPage := strings.ToLower(link) != "/wiki/main_page"
	noillegalChars := !strings.Contains(link, ":") && !strings.Contains(link, "#")
	return validPrefix && isNotMainPage && noillegalChars
}

// gets random article from metawiki API
// returns article in the form "/wiki/XXXXX"
func GetRandomArticle() (string, error) {
	req, _ := http.NewRequest("GET", metawikiEndpoint, nil)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		logErr("Could not get new article from metawiki server: %v", err)
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logErr("Could not read response from metawiki server: %v", err)
		return "", err
	}
	rArticle := &RArticleResp{}
	err = json.Unmarshal(body, &rArticle)
	if err != nil {
		logErr("could not unmarshal response from metawiki server: %v \n body: %s", err, string(body))
		return "", err
	}
	// etract response
	for _, p := range rArticle.Query.Pages {
		// return on first article
		return baseEndpoint + prefix + p.Title, nil
	}
	return "", fmt.Errorf("Could not find article in metawiki response:  %v+", rArticle)
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
	// trim current node if needed
	currentNode = CleanUrl(currentNode)
	neighborsAdded = []string{}
	// get IDs from page keys
	for i, n := range neighborNodes {
		neighborNodes[i] = CleanUrl(n)
	}
	twoWayResp, err := db.GetArticleIds(append(neighborNodes, currentNode))
	if err != nil {
		logErr("Could not get neighbor Ids %v", err)
		return neighborsAdded, err
	}
	// log out errors, if any
	for _, e := range twoWayResp.Errors {
		logErr("Error getting article ID: %s", e)
	}
	// map string => id (int)
	currentNodeId := -1
	neighborNodesIds := []int{}
	for _, entry := range twoWayResp.Entries {
		if entry.Key == currentNode {
			currentNodeId = entry.Value
		} else {
			neighborNodesIds = append(neighborNodesIds, entry.Value)
		}
	}
	// current cannot be -1
	if currentNodeId == -1 {
		logErr("Could not find reverse string => int lookup from \n resp: %v, \n currentNode: %s, \n neighbors : %v", twoWayResp.Entries, currentNode, neighborNodes)
		return neighborsAdded, errors.New("Could not find node on reverse lookup")
	}
	// post IDs to graph db
	graphResp, err := db.AddNeighbors(currentNodeId, neighborNodesIds)
	if err != nil {
		logErr("Could not POST to graph DB")
		return neighborsAdded, err
	}
	// map id => string
	for _, entry := range twoWayResp.Entries {
		for _, nAdded := range graphResp.NeighborsAdded {
			if entry.Value == nAdded {
				// add back in prefix
				neighborsAdded = append(neighborsAdded, baseEndpoint+prefix+entry.Key)
			}
		}
	}
	return neighborsAdded, err
}
