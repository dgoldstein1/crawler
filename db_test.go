package crawler

import (
	"testing"
	"os"
	"github.com/jarcoal/httpmock"
	"errors"
	"net/http"
	"encoding/json"
)

var dbEndpoint = "http://localhost:17474"
var notFoundError = `dial tcp 127.0.0.1:17474: connect: connection refused`

func TestAddToDb(t *testing.T) {
	t.Run("fails when no server found", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
		// first test bad response
		alreadyInDB, err := addEdgeIfDoesNotExist("testNode", "3")
		AssertErrorEqual(t, err, errors.New("Get http://localhost:17474/neighbors?node=testNode: " + notFoundError))
		AssertEqual(t, alreadyInDB, false)
	})
	t.Run("node already exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/neighbors?node=2",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{"5","3","6"})
			},
		)
		alreadyInDB, err := addEdgeIfDoesNotExist("2", "6")
		AssertErrorEqual(t, err, nil)
		AssertEqual(t, alreadyInDB, true)
	})
	t.Run("adds node succesfully", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/neighbors?node=2",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(404, map[string]interface{}{
					"code" : 404,
					"error" : "Node was not found",
				})
			},
		)
		httpmock.RegisterResponder("POST", dbEndpoint + "/neighbors?node=2",
			func(req *http.Request) (*http.Response, error) {
				body := make(map[string][]string)
				err := json.NewDecoder(req.Body).Decode(&body);
				if err != nil {
					t.Error(err)
				}
				return httpmock.NewJsonResponse(200, body)
			},
		)
		alreadyInDB, err := addEdgeIfDoesNotExist("2", "6")
		AssertErrorEqual(t, err, nil)
		AssertEqual(t, alreadyInDB, false)
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

		err := connectToDB()
		AssertErrorEqual(t, err, nil)
	})
}
