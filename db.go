package crawler

import (
	"os"
	"net/http"
	"io/ioutil"
)

// adds edge to DB, returns (true) if neighbor node exists
func addToDB(currentNode string, neighborNode []string) (bool, error) {
	// check to see if node already exists

	// POST node to DB

	return true, nil
}

// connects to given databse
func connectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT") + "/metrics")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	// handling error and doing stuff with body that needs to be unit tested
	return err
}
