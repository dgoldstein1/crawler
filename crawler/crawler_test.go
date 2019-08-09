package crawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	// "strings"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	// return 'true' if edge already exists
	nodesAdded := []string{}
	addEdge := func(currNode string, neighborNodes []string) ([]string, error) {
		temp := []string{}
		for _, v := range neighborNodes {
			temp = append(temp, "https://en.wikipedia.org"+v)
		}
		nodesAdded = append(nodesAdded, temp...)
		return temp, nil
	}
	// check if valid url string for crawling
	isValidCrawlLink := func(url string) bool { return true }
	// node2 := "/wiki/node2"
	newNodeRetrieved := false
	getNewNode := func() (string, error) {
		newNodeRetrieved = true
		return "https://en.wikipedia.org/wiki/String_cheese", nil
	}
	// mock fatalf
	originLogFatalf := logFatal
	defer func() { logFatal = originLogFatalf }()
	logs := []string{}
	logFatal = func(format string, args ...interface{}) {
		if len(args) > 0 {
			logs = append(logs, fmt.Sprintf(format, args))
		} else {
			logs = append(logs, format)
		}
	}

	type Test struct {
		Name             string
		StartingEndpoint string
		ConnectToDB      func() error
		MaxNodes         int32
		MinNodesAdded    int
		MaxNodesAdded    int
		NewNodeRetrieved bool
	}

	testTable := []Test{
		Test{
			Name:             "starts with endpoint if one is passed",
			StartingEndpoint: "https://en.wikipedia.org/wiki/String_cheese",
			ConnectToDB:      func() error { return nil },
			MaxNodes:         100,
			MinNodesAdded:    1,
			MaxNodesAdded:    1000,
			NewNodeRetrieved: false,
		},
		Test{
			Name:             "cannot connect to DB",
			StartingEndpoint: "https://en.wikipedia.org/wiki/String_cheese",
			ConnectToDB:      func() error { return errors.New("test error") },
			MaxNodes:         1,
			MinNodesAdded:    0,
			MaxNodesAdded:    0,
			NewNodeRetrieved: false,
		},
		Test{
			Name:             "fetches new node if endpoint is undefined",
			StartingEndpoint: "",
			ConnectToDB:      func() error { return nil },
			MaxNodes:         100,
			MinNodesAdded:    1,
			MaxNodesAdded:    1000,
			NewNodeRetrieved: true,
		},
	}

	for _, test := range testTable {
		// run tests
		t.Run(test.Name, func(t *testing.T) {
			// reset everything
			newNodeRetrieved = false
			nodesAdded = []string{}
			logs = []string{}
			Run(
				test.StartingEndpoint,
				test.MaxNodes,
				isValidCrawlLink,
				test.ConnectToDB,
				addEdge,
				getNewNode,
			)
			// make assertions
			if test.Name != "cannot connect to DB" {
				assert.True(t, test.MinNodesAdded <= len(nodesAdded))
				assert.True(t, test.MaxNodesAdded >= len(nodesAdded))
				assert.Equal(t, test.NewNodeRetrieved, newNodeRetrieved)
				assert.Equal(t, 0, len(logs))
			} else {
				assert.Equal(t, 1, len(logs))
			}
		})
	}
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
// 	endpoint := "https://en.wikipedia.org/wiki/String_cheese"
//
// 	// mock out log.Fatalf
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
// 		Crawl("https://en.wikipedia.org/wiki/String_cheese", 2, isValidCrawlLink, addEdges)
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
// func(currNode string, neighborNodes []string) ([]string, error) {
// 	temp := []string{}
// 	for _, v := range neighborNodes {
// 		temp = append(temp, "https://en.wikipedia.org"+v)
// 	}
// 	nodesAdded = append(nodesAdded, temp...)
// 	return temp, nil
// },
// 		)
//
// 		assert.Equal(t, "starting at ["+endpoint+"]", logs[0])
// 		// only add first recursion nodes, ~30,000 on second recursion
// 		assert.Equal(t, len(nodesAdded) >= 50, true)
// 		assert.Equal(t, len(nodesAdded) <= 30000, true)
// 		assert.Equal(t, []string{}, errors)
// 	})
// }
