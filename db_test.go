package crawler

import (
	"testing"
	"os"
	"github.com/jarcoal/httpmock"
	"errors"
)

var dbEndpoint = "http://localhost:17474"
var notFoundError = errors.New(`Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused`)

func TestAddToDb(t *testing.T) {
	t.Run("fails when no server found", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
		// first test bad response
		hasNeighbors, err := addToDB("testNode", []string{})
		AssertErrorEqual(t, err, notFoundError)
		AssertEqual(t, hasNeighbors, false)
		AssertErrorEqual(t, err, errors.New(`Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused`))
	})
	// mock out http endpoint
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()
	// httpmock.RegisterResponder("GET", dbEndpoint + "/neighbors?node=2",
	// 	httpmock.NewStringResponder(200, `[17,33,89,95]`))
	// // Use Client & URL from our local test server
	// err = connectToDB("2", []string{})
	// AssertErrorEqual(t, err, nil)
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	t.Run("fails when db not found", func(t *testing.T) {
		err := connectToDB()
		AssertErrorEqual(t, err, notFoundError)
	})
	t.Run("succeed when server exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/metrics",
			httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
		// Use Client & URL from our local test server
		err := connectToDB()
		AssertErrorEqual(t, err, nil)
	})
}
