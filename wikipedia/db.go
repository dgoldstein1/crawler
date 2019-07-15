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
)

var prefex = "/wiki/"
var baseEndpoint = "https://en.wikipedia.org"
func IsValidCrawlLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

// adds edge to DB, returns (true) if neighbor already in DB
func AddEdgeIfDoesNotExist(currentNode string, neighborNode string) (bool, error) {
	// get wiki IDs
	currentNodeId, err := getArticleId(strings.TrimPrefix(currentNode, baseEndpoint))
	if err != nil {
		return false, err
	}
	neighborNodeId, err := getArticleId(neighborNode)
	if err != nil {
		return false, err
	}
	// check to see if node already exists
	url := os.Getenv("GRAPH_DB_ENDPOINT") + "/neighbors"
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("node", strconv.Itoa(currentNodeId))
	req.URL.RawQuery = q.Encode()
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	// current node exists
	if resp.StatusCode != 404 {
		// check that neighbor node is not in response
		defer resp.Body.Close()
		var neighbors []int
		err = json.NewDecoder(resp.Body).Decode(&neighbors)
		if err != nil {
			return false, err
		}
		// check if neighbor is alredy added
		for _, v := range neighbors {
			if v == neighborNodeId {
				return true, nil
			}
		}
	}

	// POST node to DB
	jsonValue, _ := json.Marshal(map[string][]int{
		"neighbors": []int{neighborNodeId},
	})
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = q.Encode()
	// return the result of the POST request
	_, err = client.Do(req)
	return false, err
}


type PropertiesResponse struct {
	Parse PropertiesValues `json:"parse"`
}
type PropertiesValues struct {
	Pageid int `json:"pageid"`
	// drop title and properties keys
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
