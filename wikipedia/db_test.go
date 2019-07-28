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
var twoWayEndpoint = "http://localhost:17475"

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

func TestAddEdgesIfDoNotExist(t *testing.T) {
	os.Setenv("TWO_WAY_KV_ENDPOINT", twoWayEndpoint)
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type Test struct {
		Name             string
		Setup            func()
		CurrNode         string
		NeighborNodes    []string
		ExpectedResponse []string
		ExpectedError    error
	}
	testTable := []Test{
		Test{
			Name: "adds all neighbor nodes sucesfully",
			Setup: func() {
				// mock out DB call
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{2, 3, 4}})
					},
				)
				// mock out metadata call
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"errors": []string{"test"},
							"entries": []TwoWayEntry{
								TwoWayEntry{"/wiki/test", 1},
								TwoWayEntry{"/wiki/test1", 2},
								TwoWayEntry{"/wiki/test2", 3},
								TwoWayEntry{"/wiki/test3", 4},
							},
						})
					},
				)

			},
			CurrNode:         "/wiki/test",
			NeighborNodes:    []string{"/wiki/test1", "/wiki/test2", "/wiki/test3"},
			ExpectedResponse: []string{baseEndpoint + "/wiki/test1", baseEndpoint + "/wiki/test2", baseEndpoint + "/wiki/test3"},
			ExpectedError:    nil,
		},
		Test{
			Name: "adds all neighbor nodes sucesfully with full (non-trimmed) link",
			Setup: func() {
				// mock out DB call
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{2, 3, 4}})
					},
				)
				// mock out metadata call
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"errors": []string{"test"},
							"entries": []TwoWayEntry{
								TwoWayEntry{"/wiki/test", 1},
								TwoWayEntry{"/wiki/test1", 2},
								TwoWayEntry{"/wiki/test2", 3},
								TwoWayEntry{"/wiki/test3", 4},
							},
						})
					},
				)

			},
			CurrNode:         "https://en.wikipedia.org/wiki/test",
			NeighborNodes:    []string{"/wiki/test1", "/wiki/test2", "/wiki/test3"},
			ExpectedResponse: []string{baseEndpoint + "/wiki/test1", baseEndpoint + "/wiki/test2", baseEndpoint + "/wiki/test3"},
			ExpectedError:    nil,
		},
		Test{
			Name: "returns only neighbors added",
			Setup: func() {
				// mock out DB call
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{4}})
					},
				)
				// mock out metadata call
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"errors": []string{"test"},
							"entries": []TwoWayEntry{
								TwoWayEntry{"/wiki/test", 1},
								TwoWayEntry{"/wiki/test1", 2},
								TwoWayEntry{"/wiki/test2", 3},
								TwoWayEntry{"/wiki/test3", 4},
							},
						})
					},
				)

			},
			CurrNode:         "/wiki/test",
			NeighborNodes:    []string{"/wiki/test1", "/wiki/test1", "/wiki/test2", "/wiki/test3"},
			ExpectedResponse: []string{baseEndpoint + "/wiki/test3"},
			ExpectedError:    nil,
		},
		Test{
			Name: "fails on bad ID lookup",
			Setup: func() {
				// mock out DB call
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{2, 3, 4}})
					},
				)
				// mock out metadata call
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(500, map[string]interface{}{"error": "Could not connect to TWO_WAY_KV_ENDPOINT", "code": 500})
					},
				)
			},
			CurrNode:         "/wiki/test",
			NeighborNodes:    []string{"/wiki/test1", "/wiki/test1", "/wiki/test2", "/wiki/test3"},
			ExpectedResponse: []string{},
			ExpectedError:    errors.New("Could not connect to TWO_WAY_KV_ENDPOINT"),
		},

		Test{
			Name: "fails on bad GRAPH_DB_ENDPOINT request",
			Setup: func() {
				// mock out DB call
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(500, map[string]interface{}{"error": "Could not connect to TWO_WAY_KV_ENDPOINT", "code": 500})
					},
				)
				// mock out metadata call
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"errors": []string{"test"},
							"entries": []TwoWayEntry{
								TwoWayEntry{"/wiki/test", 1},
								TwoWayEntry{"/wiki/test1", 2},
								TwoWayEntry{"/wiki/test2", 3},
								TwoWayEntry{"/wiki/test3", 4},
							},
						})
					},
				)
			},
			CurrNode:         "/wiki/test",
			NeighborNodes:    []string{"/wiki/test1", "/wiki/test1", "/wiki/test2", "/wiki/test3"},
			ExpectedResponse: []string{},
			ExpectedError:    errors.New("Could not connect to TWO_WAY_KV_ENDPOINT"),
		},
		Test{
			Name: "fails on reverse lookup",
			Setup: func() {
				// mock out DB call
				httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=1",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{2, 3, 4}})
					},
				)
				// mock out metadata call
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"errors": []string{"test"},
							"entries": []TwoWayEntry{
								// TwoWayEntry{"/wiki/test", 1}, >> mock db not returning correct node
								TwoWayEntry{"/wiki/test1", 2},
								TwoWayEntry{"/wiki/test2", 3},
								TwoWayEntry{"/wiki/test3", 4},
							},
						})
					},
				)

			},
			CurrNode:         "/wiki/test",
			NeighborNodes:    []string{"/wiki/test1", "/wiki/test1", "/wiki/test2", "/wiki/test3"},
			ExpectedResponse: []string{},
			ExpectedError:    errors.New("Could not find node on reverse lookup"),
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Setup()
			resp, err := AddEdgesIfDoNotExist(test.CurrNode, test.NeighborNodes)
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
	os.Setenv("TWO_WAY_KV_ENDPOINT", twoWayEndpoint)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type Test struct {
		Name             string
		Articles         []string
		ExpectedResponse TwoWayResponse
		ExpectedError    error
		Setup            func()
	}
	testTable := []Test{
		Test{
			Name:     "returns response from API succesfully",
			Articles: []string{"/wiki/test", "/wiki/test1", "/wiki/test2"},
			ExpectedResponse: TwoWayResponse{
				Errors: []string{"test"},
				Entries: []TwoWayEntry{
					TwoWayEntry{"/wiki/test1", 2},
					TwoWayEntry{"/wiki/test2", 3},
					TwoWayEntry{"/wiki/test3", 4},
				},
			},
			ExpectedError: nil,
			Setup: func() {
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(200, map[string]interface{}{
							"errors": []string{"test"},
							"entries": []TwoWayEntry{
								TwoWayEntry{"/wiki/test1", 2},
								TwoWayEntry{"/wiki/test2", 3},
								TwoWayEntry{"/wiki/test3", 4},
							},
						})
					},
				)
			},
		},
		Test{
			Name:             "handles 500 code response",
			Articles:         []string{"/wiki/test", "/wiki/test1", "/wiki/test2"},
			ExpectedResponse: TwoWayResponse{},
			ExpectedError:    errors.New("server error"),
			Setup: func() {
				httpmock.RegisterResponder("POST", twoWayEndpoint+"/entries",
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(500, map[string]interface{}{"error": "server error", "code": 500})
					},
				)
			},
		},
		Test{
			Name:             "returns error on bad endpoint",
			Articles:         []string{"/wiki/test", "/wiki/test1", "/wiki/test2"},
			ExpectedResponse: TwoWayResponse{},
			ExpectedError:    errors.New("Post http://localhost:17475/entries: no responder found"),
			Setup:            func() {},
		},
	}

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
