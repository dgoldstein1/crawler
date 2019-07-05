package crawler

import (
	"testing"
	"os"
	"github.com/jarcoal/httpmock"
	"errors"
	"net/http"
)

var dbEndpoint = "http://localhost:17474"
var notFoundError = `dial tcp 127.0.0.1:17474: connect: connection refused`

func TestAddToDb(t *testing.T) {
	t.Run("fails when no server found", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
		// first test bad response
		hasNeighbors, err := addToDB("testNode", []string{})
		AssertErrorEqual(t, err, errors.New("Get http://localhost:17474/neighbors?neighbors=testNode: " + notFoundError))
		AssertEqual(t, hasNeighbors, false)
	})
	t.Run("node already exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/neigihbors?node=2",
			func(req *http.Request) (*http.Response, error) {
				response := []string{"15", "36"}
				return httpmock.NewJsonResponse(200, response)
			},
		)
		// Use Client & URL from our local test server
		hasNeighbor, err := addToDB("2", []string{})
		AssertErrorEqual(t, err, nil)
		AssertEqual(t, hasNeighbor, true)
	})
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	t.Run("fails when db not found", func(t *testing.T) {
		err := connectToDB()
		AssertErrorEqual(t, err, errors.New("Get http://localhost:17474/metrics: " + notFoundError))
	})
	t.Run("succeed when server exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/metrics",
			httpmock.NewStringResponder(200, `TEST`))
		// Use Client & URL from our local test server
		err := connectToDB()
		AssertErrorEqual(t, err, nil)
	})
}
