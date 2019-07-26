package wikipedia

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// globals
var logErr = log.Errorf
var prefex = "/wiki/"
var baseEndpoint = "https://en.wikipedia.org"
var c = colly.NewCollector()

// determines if is good link to crawl on
func IsValidCrawlLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

// adds edge to DB, returns new neighbors added (to crawl on)
func AddEdgesIfDoNotExist(currentNode string, neighborNodes []string) ([]string, error) {
	return []string{}, nil
}

// posts possible new edges to GRAPH_DB_ENDPOINT
func addNeighbors(curr int, neighborIds []int) (resp GraphResponseSuccess, err error) {
	// POST new neighbors to db
	jsonValue, _ := json.Marshal(map[string][]int{
		"neighbors": neighborIds,
	})
	url := os.Getenv("GRAPH_DB_ENDPOINT") + "/edges"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("node", strconv.Itoa(curr))
	req.URL.RawQuery = q.Encode()

	// return the result of the POST request
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	// assert response is 200
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return resp, err
		}
		errResp := GraphResponseError{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return resp, err
		}
		// fails with error
		return resp, errors.New(errResp.Error)
	}

	// 200 level response, continue as normal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	resp = GraphResponseSuccess{}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

// gets wikipedia int id from article url
func getArticleIds(articles []string) (resp TwoWayResponse, err error) {
	// create array of entries
	entries := []TwoWayEntry{}
	for _, s := range articles {
		entries = append(entries, TwoWayEntry{s, 0})
	}
	// post to endpoint
	return resp, err
}

// connects to given databse and initializes scraper
func ConnectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT") + "/metrics")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
