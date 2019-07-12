package wikipedia

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"fmt"
)

func IsValidCrawlLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

// adds edge to DB, returns (true) if neighbor already in DB
func AddEdgeIfDoesNotExist(currentNode string, neighborNode string) (bool, error) {
	// check to see if node already exists
	url := os.Getenv("GRAPH_DB_ENDPOINT") + "/neighbors"
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("node", currentNode)
	req.URL.RawQuery = q.Encode()
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	// 404 is current node does not exist
	if resp.StatusCode != 404 {
		// check that neighbor node is not in response
		defer resp.Body.Close()
		var neighbors []string
		err = json.NewDecoder(resp.Body).Decode(&neighbors)
		if err != nil {
			return false, err
		}
		for _, v := range neighbors {
			// neighbor node found for this node
			if v == neighborNode {
				return true, nil
			}
		}
	}
	// no neighbor node, POST node to DB
	jsonValue, _ := json.Marshal(map[string][]string{
		"neighbors": []string{neighborNode},
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
	url := "https://en.wikipedia.org/w/api.php"
	req, _ := http.NewRequest("GET", url, nil)
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
	if (err == nil && props.Parse.Pageid == 0) {
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
