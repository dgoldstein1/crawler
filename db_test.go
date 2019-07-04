package crawler

import (
	"testing"
	"os"
	"github.com/jarcoal/httpmock"
	"errors"
)

var dbEndpoint = "http://localhost:17474"
var notFoundError = `Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused`

func TestAddToDb(t *testing.T) {
	// os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	// // first test bad response
	// err := connectToDB()
	// AssertErrorEqual(t, err, errors.New(`Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused`))
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	// first test bad response
	err := connectToDB()
	AssertErrorEqual(t, err, errors.New(notFoundError))
	// mock out http endpoint
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
		// Exact URL match
	httpmock.RegisterResponder("GET", dbEndpoint + "/metrics",
		httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
	// Use Client & URL from our local test server
	err = connectToDB()
	AssertErrorEqual(t, err, nil)
}
