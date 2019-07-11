package wikipedia

import (
	"testing"
	"os"
	"github.com/jarcoal/httpmock"
	"net/http"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var dbEndpoint = "http://localhost:17474"
var notFoundError = `dial tcp 127.0.0.1:17474: connect: connection refused`

func TestIsValidCrawlLink(t *testing.T) {
  t.Run("does not crawl on links with ':'", func(t *testing.T) {
    assert.Equal(t, IsValidCrawlLink("/wiki/Category:Spinash"), false)
    assert.Equal(t, IsValidCrawlLink("/wiki/Test:"), false)
  })
  t.Run("does not crawl on links not starting with '/wiki/'", func(t *testing.T ){
    assert.Equal(t, IsValidCrawlLink("https://wikipedia.org"), false)
    assert.Equal(t, IsValidCrawlLink("/wiki"), false)
    assert.Equal(t, IsValidCrawlLink("wikipedia/wiki/"), false)
    assert.Equal(t, IsValidCrawlLink("/wiki/binary"), true)
  })
}

func TestAddToDb(t *testing.T) {
	t.Run("fails when no server found", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
		// first test bad response
		alreadyInDB, err := AddEdgeIfDoesNotExist("testNode", "3")
		assert.EqualError(t, err, "Get http://localhost:17474/neighbors?node=testNode: " + notFoundError)
		assert.Equal(t, alreadyInDB, false)
	})
	t.Run("neighbor node already exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/neighbors?node=2",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{"5","3","6"})
			},
		)
		alreadyInDB, err := AddEdgeIfDoesNotExist("2", "6")
		assert.Nil(t, err)
		assert.Equal(t, alreadyInDB, true)
	})
	t.Run("adds node when current node doesnt exist (404)", func(t *testing.T){
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
		alreadyInDB, err := AddEdgeIfDoesNotExist("2", "6")
		assert.Nil(t, err)
		assert.Equal(t, alreadyInDB, false)
	})
	t.Run("adds node when current exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/neighbors?node=2",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, []string{"5","3","7"})
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
		alreadyInDB, err := AddEdgeIfDoesNotExist("2", "6")
		assert.Nil(t, err)
		assert.Equal(t, alreadyInDB, false)
	})
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	t.Run("fails when db not found", func(t *testing.T) {
		err := ConnectToDB()
		assert.EqualError(t, err, "Get http://localhost:17474/metrics: " + notFoundError)
	})
	t.Run("succeed when server exists", func(t *testing.T){
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint + "/metrics",
			httpmock.NewStringResponder(200, `TEST`))

		err := ConnectToDB()
		assert.Nil(t, err)
	})
}
