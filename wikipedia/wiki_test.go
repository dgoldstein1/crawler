package wikipedia

import (
	"errors"
	"fmt"
	db "github.com/dgoldstein1/crawler/db"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
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

func TestGetRandomArticle(t *testing.T) {
	errorsLogged := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errorsLogged = append(errorsLogged, fmt.Sprintf(format, args))
		} else {
			errorsLogged = append(errorsLogged, format)
		}
	}

	type Test struct {
		Name             string
		MockedRequest    string
		ExpectedResponse string
		ExpectedError    string
	}

	testTable := []Test{
		Test{
			Name:             "succesful",
			MockedRequest:    `{"batchcomplete":"","continue":{"grncontinue":"0.369259750651|0.369260921533|12247122|0","continue":"grncontinue||"},"query":{"pages":{"9820486":{"pageid":9820486,"ns":0,"title":"Oregon Bicycle Racing Association","extract":"The Oregon Bicycle Racing Association is a bicycle racing organization based in the U.S. state of Oregon."}}},"limits":{"extracts":20}}`,
			ExpectedResponse: "https://en.wikipedia.org/wiki/Oregon Bicycle Racing Association",
			ExpectedError:    "",
		},
		Test{
			Name:             "bad json",
			MockedRequest:    `XXXXXXbatchcomplete":"","continue":{"grncontinue":"0.369259750651|0.369260921533|12247122|0","continue":"grncontinue||"},"query":{"pages":{"9820486":{"pageid":9820486,"ns":0,"title":"Oregon Bicycle Racing Association","extract":"The Oregon Bicycle Racing Association is a bicycle racing organization based in the U.S. state of Oregon."}}},"limits":{"extracts":20}}`,
			ExpectedResponse: "",
			ExpectedError:    "invalid character 'X' looking for beginning of value",
		},
		Test{
			Name:             "ENDPOINT_NOT_FOUND",
			MockedRequest:    "",
			ExpectedResponse: "",
			ExpectedError:    "Get http://BAD_ENDPOINT: dial tcp: lookup BAD_ENDPOINT",
		},
	}

	tempEndpointVar := ""
	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// mock out endpoint
			if test.Name == "ENDPOINT_NOT_FOUND" {
				tempEndpointVar = metawikiEndpoint
				metawikiEndpoint = "http://BAD_ENDPOINT"
				timeout = time.Duration(1 * time.Second)
			} else {
				httpmock.Activate()
				httpmock.RegisterResponder("GET", metawikiEndpoint,
					httpmock.NewStringResponder(200, test.MockedRequest))
			}
			// run test
			a, err := GetRandomArticle()
			assert.Equal(t, test.ExpectedResponse, a)
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), test.ExpectedError))
				assert.Equal(t, 1, len(errorsLogged))
			} else {
				assert.Equal(t, "", test.ExpectedError)
				assert.Equal(t, 0, len(errorsLogged))
			}
			// reset
			httpmock.DeactivateAndReset()
			errorsLogged = []string{}
			timeout = time.Duration(5 * time.Second)
			metawikiEndpoint = tempEndpointVar
		})
	}

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
							"entries": []db.TwoWayEntry{
								db.TwoWayEntry{"/wiki/test", 1},
								db.TwoWayEntry{"/wiki/test1", 2},
								db.TwoWayEntry{"/wiki/test2", 3},
								db.TwoWayEntry{"/wiki/test3", 4},
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
							"entries": []db.TwoWayEntry{
								db.TwoWayEntry{"/wiki/test", 1},
								db.TwoWayEntry{"/wiki/test1", 2},
								db.TwoWayEntry{"/wiki/test2", 3},
								db.TwoWayEntry{"/wiki/test3", 4},
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
							"entries": []db.TwoWayEntry{
								db.TwoWayEntry{"/wiki/test", 1},
								db.TwoWayEntry{"/wiki/test1", 2},
								db.TwoWayEntry{"/wiki/test2", 3},
								db.TwoWayEntry{"/wiki/test3", 4},
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
							"entries": []db.TwoWayEntry{
								db.TwoWayEntry{"/wiki/test", 1},
								db.TwoWayEntry{"/wiki/test1", 2},
								db.TwoWayEntry{"/wiki/test2", 3},
								db.TwoWayEntry{"/wiki/test3", 4},
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
							"entries": []db.TwoWayEntry{
								// db.TwoWayEntry{"/wiki/test", 1}, >> mock db not returning correct node
								db.TwoWayEntry{"/wiki/test1", 2},
								db.TwoWayEntry{"/wiki/test2", 3},
								db.TwoWayEntry{"/wiki/test3", 4},
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
