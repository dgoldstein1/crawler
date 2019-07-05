package crawler

import (
	"os"
	"net/http"
	"io/ioutil"
	// "encoding/json"
	"time"
)

// adds edge to DB, returns (true) if already in DB
func addToDB(currentNode string, neighborNode []string) (bool, error) {
	// check to see if node already exists
	req, _ := http.NewRequest("GET", os.Getenv("GRAPH_DB_ENDPOINT") + "/neighbors", nil)
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
	// assert that node does not exist
	if (resp.StatusCode != 404) {
		return true, nil
	}

	// POST node to DB

	return false, nil
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
