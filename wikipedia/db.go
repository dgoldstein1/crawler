package wikipedia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"log"
	"errors"
)

// globals
var logMsg = log.Printf
var prefex = "/wiki/"
var baseEndpoint = "https://en.wikipedia.org"

// determines if is good link to crawl on
func IsValidCrawlLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

// adds edge to DB, returns new neighbors added (to crawl on)
func AddEdgesIfDoNotExist(currentNode string, neighborNodes []string) ([]string, error) {
	// get wiki IDs
	currentNodeId, err := getArticleId(strings.TrimPrefix(currentNode, baseEndpoint))
	if err != nil {
		return []string{}, err
	}
	// make a map of id[value]
	neighborsMap := make(map[int]string)
	neighborsIds := []int{}
	for _, n := range neighborNodes {
		neighborNodeId, err := getArticleId(n)
		if err != nil {
			logMsg("ERROR: %s", err.Error())
		} else {
			neighborsMap[neighborNodeId] = n
			neighborsIds = append(neighborsIds, neighborNodeId)
		}
	}

	// POST new neighbors to db
	jsonValue, _ := json.Marshal(map[string][]int{
		"neighbors": neighborsIds,
	})
	url := os.Getenv("GRAPH_DB_ENDPOINT") + "/edges"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("node", strconv.Itoa(currentNodeId))
	req.URL.RawQuery = q.Encode()

	// return the result of the POST request
   client := http.Client{
     Timeout: time.Duration(5 * time.Second),
  }
	res, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	// assert response is 200
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []string{}, err
		}
		errResp := GraphResponseError{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return []string{}, err
		}
		// fails with error
		return []string{}, errors.New(errResp.Error)
	}

	// 200 level response, continue as normal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}
	resp := GraphResponseSuccess{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return []string{}, err
	}
	newEdgesNodes := resp.NeighborsAdded
	// compare new ids to
	nodesAdded := []string{}
	for _, n := range newEdgesNodes {
		nodesAdded = append(nodesAdded, neighborsMap[n])
	}
	return nodesAdded, nil
}


// gets wikipedia int id from article url
func getArticleId(page string) (int, error) {
	parsedPage := strings.TrimPrefix(page, "/wiki/")
	req, _ := http.NewRequest("GET", os.Getenv("WIKI_API_ENDPOINT"), nil)
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("action", "parse")
	q.Add("page", parsedPage)
	q.Add("prop", "properties")
	req.URL.RawQuery = q.Encode()
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	props := &PropertiesResponse{}
	err = json.Unmarshal(body, &props)
	if err == nil && props.Parse.Pageid == 0 {
		err = fmt.Errorf("Page not found '%s'", page)
	}
	return props.Parse.Pageid, err
}

// connects to given databse
func ConnectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT") + "/metrics")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
