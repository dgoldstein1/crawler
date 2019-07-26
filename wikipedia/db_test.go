package wikipedia

import (
	// "fmt"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var dbEndpoint = "http://localhost:17474"

func TestIsValidCrawlLink(t *testing.T) {
	t.Run("does not crawl on links with ':'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/wiki/Category:Spinash"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/Test:"), false)
	})
	t.Run("does not crawl on links not starting with '/wiki/'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://wikipedia.org"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki"), false)
		assert.Equal(t, IsValidCrawlLink("wikipedia/wiki/"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/binary"), true)
	})
}

func TestAddNeighbors(t *testing.T) {
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type Test struct {
		Name             string
		Setup            func()
		CurrNode         int
		NeighborIds      []int
		ExpectedResponse GraphResponseSuccess
		ExpectedError    error
	}
	testTable := []Test{
		Test{
			Name: "returns correct response",
			Setup: func() {
				// Exact URL match
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{2, 3, 4}})
					},
				)
			},
			CurrNode:    1,
			NeighborIds: []int{2, 3, 4},
			ExpectedResponse: GraphResponseSuccess{
				NeighborsAdded: []int{2, 3, 4},
			},
			ExpectedError: nil,
		},
		Test{
			Name: "returns error on 500 level code",
			Setup: func() {
				// Exact URL match
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(500, map[string]interface{}{"error": "Not Found", "code": 500})
					},
				)
			},
			CurrNode:         1,
			NeighborIds:      []int{2, 3, 4},
			ExpectedResponse: GraphResponseSuccess{},
			ExpectedError:    errors.New("Not Found"),
		},
		Test{
			Name:             "Bad endpoint",
			Setup:            func() {},
			CurrNode:         1,
			NeighborIds:      []int{2, 3, 4},
			ExpectedResponse: GraphResponseSuccess{},
			ExpectedError:    errors.New("Post http://localhost:17474/edges?node=1: no responder found"),
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Setup()
			resp, err := addNeighbors(test.CurrNode, test.NeighborIds)
			if err != nil && test.ExpectedError != nil {
				assert.Equal(t, test.ExpectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.ExpectedError, err)
			}
			assert.Equal(t, test.ExpectedResponse, resp)
			httpmock.Reset()
		})
	}

}

func TestGetArticleIds(t *testing.T) {
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type Test struct {
		Name             string
		Setup            func()
		ExpectedResponse GraphResponseSuccess
		ExpectedError    error
		Articles         []string
	}
	testTable := []Test{}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Setup()
			resp, err := getArticleIds(test.Articles)
			if err != nil && test.ExpectedError != nil {
				assert.Equal(t, test.ExpectedError.Error(), err.Error())
			} else {
				assert.Equal(t, test.ExpectedError, err)
			}
			assert.Equal(t, test.ExpectedResponse, resp)
			httpmock.Reset()
		})
	}
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	t.Run("fails when db not found", func(t *testing.T) {
		err := ConnectToDB()
		assert.EqualError(t, err, "Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused")
	})
	t.Run("succeed when server exists", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint+"/metrics",
			httpmock.NewStringResponder(200, `TEST`))

		err := ConnectToDB()
		assert.Nil(t, err)
	})
}
