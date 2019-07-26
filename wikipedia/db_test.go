package wikipedia

import (
	// "fmt"
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
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Setup()
			resp, err := addNeighbors(test.CurrNode, test.NeighborIds)
			assert.Equal(t, test.ExpectedError, err)
			assert.Equal(t, test.ExpectedResponse, resp)
			httpmock.Reset()
		})
	}

}

// func TestAddToDb(t *testing.T) {
// 	// keep errors in array
// 	errors := []string{}
// 	logErr = func(format string, args ...interface{}) {
// 		if len(args) > 0 {
// 			errors = append(errors, fmt.Sprintf(format, args))
// 		} else {
// 			errors = append(errors, format)
// 		}
// 	}
// 	t.Run("fails when no server found", func(t *testing.T) {
// 		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
// 		// first test bad response
// 		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet", []string{"/wiki/Animal"})
// 		assert.EqualError(t, err, "Post http://localhost:17474/edges?node=25079: dial tcp 127.0.0.1:17474: connect: connection refused")
// 		assert.Equal(t, []string{}, newNodes)
// 	})
// 	t.Run("returns error when current node doesnt exist (404)", func(t *testing.T) {
// 		// mock out http endpoint
// 		httpmock.Activate()
// 		defer httpmock.DeactivateAndReset()
// 		// Exact URL match
// 		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
// 			func(req *http.Request) (*http.Response, error) {
// 				return httpmock.NewJsonResponse(404, map[string]interface{}{
// 					"code":  404,
// 					"error": "Node was not found",
// 				})
// 			},
// 		)
//
// 		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal"})
// 		assert.EqualError(t, err, "Node was not found")
// 		assert.Equal(t, newNodes, []string{})
// 	})
// 	t.Run("succesfully adds neighbor nodes", func(t *testing.T) {
// 		errors = []string{}
// 		// mock out http endpoint
// 		httpmock.Activate()
// 		defer httpmock.DeactivateAndReset()
// 		// Exact URL match
//
// 		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
// 			func(req *http.Request) (*http.Response, error) {
// 				return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{11039790}})
// 			},
// 		)
//
// 		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal"})
// 		assert.Nil(t, err)
// 		assert.Equal(t, errors, []string{})
// 		assert.Equal(t, newNodes, []string{"https://en.wikipedia.org/wiki/Animal"})
// 	})
// 	t.Run("only returns new neighbors", func(t *testing.T) {
// 		// mock out http endpoint
// 		// mock out http endpoint
// 		httpmock.Activate()
// 		defer httpmock.DeactivateAndReset()
// 		// Exact URL match
//
// 		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
// 			func(req *http.Request) (*http.Response, error) {
// 				return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{11039790}})
// 			},
// 		)
//
// 		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal", "/wiki/Petula_clark"})
// 		assert.Nil(t, err)
// 		assert.Equal(t, newNodes, []string{"https://en.wikipedia.org/wiki/Animal"})
// 	})
// }

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
