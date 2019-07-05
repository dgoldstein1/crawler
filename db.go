package crawler

import (
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"bytes"
)

// adds edge to DB, returns (true) if neighbor already in DB
func addEdgeIfDoesNotExist(currentNode string, neighborNode string) (bool, error) {
	// check to see if node already exists
	url := os.Getenv("GRAPH_DB_ENDPOINT") + "/neighbors"
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("node", currentNode)
	req.URL.RawQuery = q.Encode()
	client := http.Client{
		Timeout : time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	// 404 is current node does not exist
	if (resp.StatusCode != 404) {
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
	jsonValue, _ := json.Marshal( map[string][]string {
		"neighbors" : []string{neighborNode},
	})
  req, _ = http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
  req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = q.Encode()
	// return the result of the POST request
	_, err = client.Do(req)
	return false, err
}

// connects to given databse
func connectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT") + "/metrics")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
