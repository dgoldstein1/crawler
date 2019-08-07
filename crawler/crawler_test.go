package crawler

import (
	// "fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	// "strings"
	"testing"
)

func TestRun(t *testing.T) {
	// return 'true' if edge already exists
	edgesAdded := make(map[string][]string)
	addEdge := func(node string, neighbors []string) ([]string, error) {
		edgesAdded[node] = neighbors
		return neighbors, nil
	}
	// establishes initial connection to DB
	connectToDB := func() error {
		return nil
	}
	// check if valid url string for crawling
	isValidCrawlLink := func(url string) bool {
		return true
	}
	// retrieves new node if current expires
	node1 := "/wiki/node1"
	// node2 := "/wiki/node2"
	getNewNode := func() (string, error) {
		return node1, nil
	}

	// mock configurations
	wikiEndpoint := "http://wiki-endpoint"
	nodesScraped := []string{}
	setup := func() {
		httpmock.Activate()
		httpmock.RegisterResponder("GET", wikiEndpoint+node1,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(500, map[string]interface{}{"error": "server error", "code": 500})
			},
		)
	}
	tearDown := func() {
		edgesAdded = make(map[string][]string)
		httpmock.Reset()
		nodesScraped = []string{}
	}

	t.Run("starts with endpoint if one is passed", func(t *testing.T) {
		setup()
		Run(
			"/wiki/test1",
			1,
			isValidCrawlLink,
			connectToDB,
			addEdge,
			getNewNode,
		)
		assert.Equal(t, []string{}, nodesScraped)
		tearDown()

	})
	t.Run("starts with random node if none is passed", func(t *testing.T) {
		// TODO
	})
	t.Run("runs until max approximate nodes is reached", func(t *testing.T) {
		// TODO:
	})
	t.Run("runs until max retries on get node is reached", func(t *testing.T) {
		// TODO
	})
}

//
// func TestCrawl(t *testing.T) {
// 	isValidCrawlLink := func(url string) bool {
// 		return strings.HasPrefix(url, "/wiki/") && !strings.Contains(url, ":")
// 	}
// 	nodesAdded := []string{}
// 	addEdges := func(currNode string, neighborNodes []string) ([]string, error) {
// 		nodesAdded = append(nodesAdded, neighborNodes...)
// 		return neighborNodes, nil
// 	}
// 	connectToDB := func() error { return nil }
// 	endpoint := "https://en.wikipedia.org/wiki/String_cheese"
//
// 	// mock out log.Fatalf
// 	originLogPrintf := logMsg
// 	defer func() { logMsg = originLogPrintf }()
// 	logs := []string{}
// 	logMsg = func(format string, args ...interface{}) {
// 		if len(args) > 0 {
// 			logs = append(logs, fmt.Sprintf(format, args))
// 		} else {
// 			logs = append(logs, format)
// 		}
// 	}
//
// 	// mute warnings
// 	originLogWarn := logWarn
// 	defer func() { logWarn = originLogWarn }()
// 	logWarn = func(format string, args ...interface{}) {}
//
// 	// keep errors in array
// 	errors := []string{}
// 	logErr = func(format string, args ...interface{}) {
// 		if len(args) > 0 {
// 			errors = append(errors, fmt.Sprintf(format, args))
// 		} else {
// 			errors = append(errors, format)
// 		}
// 	}
//
// 	t.Run("works with isValidCrawlLink", func(t *testing.T) {
// 		nodesAdded = []string{}
// 		// function doing setup of tests
// 		Crawl("https://en.wikipedia.org/wiki/String_cheese", 2, isValidCrawlLink, connectToDB, addEdges)
// 		t.Run("only filters on links starting with regex", func(t *testing.T) {
// 			errors = []string{}
// 			for _, url := range nodesAdded {
// 				assert.Equal(t, strings.HasPrefix(url, "/wiki/"), true)
// 			}
// 			assert.Equal(t, []string{}, errors)
// 		})
// 		t.Run("only filters on links starting with regex", func(t *testing.T) {
// 			errors = []string{}
// 			for _, url := range nodesAdded {
// 				if strings.Contains(url, ":") {
// 					t.Errorf("Did not expect '%s' to contain ':'", url)
// 				}
// 			}
// 			assert.Equal(t, []string{}, errors)
// 		})
// 	})
//
// 	t.Run("adds nodes correctly", func(t *testing.T) {
// 		nodesAdded = []string{}
// 		logs = []string{}
// 		errors = []string{}
// 		Crawl(
// 			endpoint,
// 			100,
// 			isValidCrawlLink,
// 			connectToDB,
// 			func(currNode string, neighborNodes []string) ([]string, error) {
// 				temp := []string{}
// 				for _, v := range neighborNodes {
// 					temp = append(temp, "https://en.wikipedia.org"+v)
// 				}
// 				nodesAdded = append(nodesAdded, temp...)
// 				return temp, nil
// 			},
// 		)
//
// 		assert.Equal(t, "starting at ["+endpoint+"]", logs[0])
// 		// only add first recursion nodes, ~30,000 on second recursion
// 		assert.Equal(t, len(nodesAdded) >= 50, true)
// 		assert.Equal(t, len(nodesAdded) <= 30000, true)
// 		assert.Equal(t, []string{}, errors)
// 	})
// }
