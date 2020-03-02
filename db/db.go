package db

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var logErr = log.Errorf
var timeout = time.Duration(5 * time.Second)

// posts possible new edges to GRAPH_DB_ENDPOINT
func AddNeighbors(curr int, neighborIds []int) (resp GraphResponseSuccess, err error) {
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
		Timeout: timeout,
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
func GetArticleIds(articles []string) (resp TwoWayResponse, err error) {
	// create array of entries
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(articles)
	// post to endpoint
	url := os.Getenv("TWO_WAY_KV_ENDPOINT") + "/entries"
	req, _ := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("muteAlreadyExistsError", "true")
	req.URL.RawQuery = q.Encode()
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	// read out response
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
	// succesful request
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	resp = TwoWayResponse{}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

// connects to given databse and initializes scraper
func ConnectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

// adds edge to DB, returns new neighbors added (to crawl on)
func AddEdgesIfDoNotExist(
	currentNode string,
	neighborNodes []string,
	cleanUrl func(string) string,
	baseEndpoint string,
) (
	neighborsAdded []string,
	err error,
) {
	// trim current node if needed
	currentNode = cleanUrl(currentNode)
	// get IDs from page keys
	// make big map of  cleanName : originalName for later
	nodes := make(map[string]string)
	for i, n := range neighborNodes {
		neighborNodes[i] = cleanUrl(n)
		nodes[cleanUrl(n)] = n
	}
	twoWayResp, err := GetArticleIds(append(neighborNodes, currentNode))
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
	graphResp, err := AddNeighbors(currentNodeId, neighborNodesIds)
	if err != nil {
		logErr("Could not POST to graph DB")
		return neighborsAdded, err
	}
	// map id => string
	for _, entry := range twoWayResp.Entries {
		for _, nAdded := range graphResp.NeighborsAdded {
			nAddedInt, _ := strconv.Atoi(nAdded)
			if entry.Value == nAddedInt {
				// add back in prefix
				neighborsAdded = append(neighborsAdded, baseEndpoint+nodes[entry.Key])
			}
		}
	}
	return neighborsAdded, err
}
