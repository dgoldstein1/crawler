package crawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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
	connectToDB := func() error { return nil }
	// check if valid url string for crawling
	isValidCrawlLink := func(url string) bool { return true }
	// node2 := "/wiki/node2"
	getNewNode := func() (string, error) {
		return "/wiki/node1", nil
	}

	type Test struct {
		Name             string
		StartingEndpoint string
		MaxNodes         int32
		MaxtRetries      int32
		MinNodesAdded    int32
		MaxNodesAdded    int32
		MaxRetries       string
	}

	testTable := []Test{
		Test{
			Name:             "starts with endpoint if one is passed",
			StartingEndpoint: "/wiki/String_cheese",
			MaxNodes:         1,
			MaxtRetries:      0,
			MinNodesAdded:    1,
			MaxNodesAdded:    25,
			MaxRetries:       "5",
		},
	}

	for _, test := range testTable {
		// run tests
		t.Run(test.Name, func(t *testing.T) {
			// reset everything
			totalNodesAdded = asyncInt(0)
			maxDepth = asyncInt(0)
			os.Setenv("MAX_NEW_NODE_RETRIES", test.MaxRetries)
			Run(
				test.StartingEndpoint,
				test.MaxNodes,
				isValidCrawlLink,
				connectToDB,
				addEdge,
				getNewNode,
			)
			// make assertions
			assert.True(t, test.MinNodesAdded < totalNodesAdded.get())
			assert.True(t, test.MaxNodesAdded > totalNodesAdded.get())
		})
	}
}

func TestCrawl(t *testing.T) {
	isValidCrawlLink := func(url string) bool {
		return strings.HasPrefix(url, "/wiki/") && !strings.Contains(url, ":")
	}
	nodesAdded := []string{}
	addEdges := func(currNode string, neighborNodes []string) ([]string, error) {
		nodesAdded = append(nodesAdded, neighborNodes...)
		return neighborNodes, nil
	}
	endpoint := "https://en.wikipedia.org/wiki/String_cheese"

	// mock out log.Fatalf
	originLogPrintf := logMsg
	defer func() { logMsg = originLogPrintf }()
	logs := []string{}
	logMsg = func(format string, args ...interface{}) {
		if len(args) > 0 {
			logs = append(logs, fmt.Sprintf(format, args))
		} else {
			logs = append(logs, format)
		}
	}

	// mute warnings
	originLogWarn := logWarn
	defer func() { logWarn = originLogWarn }()
	logWarn = func(format string, args ...interface{}) {}

	// keep errors in array
	errors := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	t.Run("works with isValidCrawlLink", func(t *testing.T) {
		nodesAdded = []string{}
		// function doing setup of tests
		Crawl("https://en.wikipedia.org/wiki/String_cheese", 2, isValidCrawlLink, addEdges)
		t.Run("only filters on links starting with regex", func(t *testing.T) {
			errors = []string{}
			for _, url := range nodesAdded {
				assert.Equal(t, strings.HasPrefix(url, "/wiki/"), true)
			}
			assert.Equal(t, []string{}, errors)
		})
		t.Run("only filters on links starting with regex", func(t *testing.T) {
			errors = []string{}
			for _, url := range nodesAdded {
				if strings.Contains(url, ":") {
					t.Errorf("Did not expect '%s' to contain ':'", url)
				}
			}
			assert.Equal(t, []string{}, errors)
		})
	})

	t.Run("adds nodes correctly", func(t *testing.T) {
		nodesAdded = []string{}
		logs = []string{}
		errors = []string{}
		Crawl(
			endpoint,
			100,
			isValidCrawlLink,
			func(currNode string, neighborNodes []string) ([]string, error) {
				temp := []string{}
				for _, v := range neighborNodes {
					temp = append(temp, "https://en.wikipedia.org"+v)
				}
				nodesAdded = append(nodesAdded, temp...)
				return temp, nil
			},
		)

		assert.Equal(t, "starting at ["+endpoint+"]", logs[0])
		// only add first recursion nodes, ~30,000 on second recursion
		assert.Equal(t, len(nodesAdded) >= 50, true)
		assert.Equal(t, len(nodesAdded) <= 30000, true)
		assert.Equal(t, []string{}, errors)
	})
}
