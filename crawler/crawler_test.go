package crawler

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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

	logMsg = func(format string, args ...interface{}) {}

	type Test struct {
		Name             string
		StartingEndpoint string
		ConnectToDB      func() error
		GetNewNode       func() (string, error)
		FilterPage       FilterPageFunction
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
			GetNewNode: func() (string, error) {
				newNodeRetrieved = true
				return "https://en.wikipedia.org/wiki/String_cheese", nil
			},
			FilterPage:       func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
			MaxNodes:         100,
			MinNodesAdded:    1,
			MaxNodesAdded:    1000,
			NewNodeRetrieved: false,
		},
		Test{
			Name:             "cannot connect to DB",
			StartingEndpoint: "https://en.wikipedia.org/wiki/String_cheese",
			ConnectToDB:      func() error { return errors.New("test error") },
			GetNewNode: func() (string, error) {
				newNodeRetrieved = true
				return "https://en.wikipedia.org/wiki/String_cheese", nil
			},
			FilterPage:       func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
			MaxNodes:         1,
			MinNodesAdded:    0,
			MaxNodesAdded:    0,
			NewNodeRetrieved: false,
		},
		Test{
			Name:             "fetches new node if endpoint is undefined",
			StartingEndpoint: "",
			ConnectToDB:      func() error { return nil },
			GetNewNode: func() (string, error) {
				newNodeRetrieved = true
				return "https://en.wikipedia.org/wiki/String_cheese", nil
			},
			FilterPage:       func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
			MaxNodes:         100,
			MinNodesAdded:    1,
			MaxNodesAdded:    1000,
			NewNodeRetrieved: true,
		},
		Test{
			Name:             "fails when cannot fetch new node",
			StartingEndpoint: "",
			ConnectToDB:      func() error { return nil },
			GetNewNode: func() (string, error) {
				newNodeRetrieved = true
				return "https://en.wikipedia.org/wiki/String_cheese", errors.New("Could not retrieve new node")
			},
			FilterPage:       func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
			MaxNodes:         100,
			MinNodesAdded:    1,
			MaxNodesAdded:    1000,
			NewNodeRetrieved: true,
		},
		Test{
			Name:             "logs error when filter page fails",
			StartingEndpoint: "https://en.wikipedia.org/wiki/String_cheese",
			ConnectToDB:      func() error { return nil },
			GetNewNode: func() (string, error) {
				newNodeRetrieved = true
				return "https://en.wikipedia.org/wiki/String_cheese", nil
			},
			FilterPage:       func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, errors.New("BAD FILTER PAGE") },
			MaxNodes:         100,
			MinNodesAdded:    1,
			MaxNodesAdded:    1000,
			NewNodeRetrieved: false,
		},
	}

	for _, test := range testTable {
		// run tests
		t.Run(test.Name, func(t *testing.T) {
			os.Setenv("MAX_APPROX_NODES", string(test.MaxNodes))
			defer os.Unsetenv("MAX_APPROX_NODES")
			Run(
				test.StartingEndpoint,
				isValidCrawlLink,
				test.ConnectToDB,
				addEdge,
				test.GetNewNode,
				test.FilterPage,
			)
			// make assertions
			if test.Name != "cannot connect to DB" && test.Name != "fails when cannot fetch new node" {
				assert.True(t, test.MinNodesAdded <= len(nodesAdded))
				assert.True(t, test.MaxNodesAdded >= len(nodesAdded))
				assert.Equal(t, test.NewNodeRetrieved, newNodeRetrieved)
				assert.Equal(t, 0, len(logs))
			} else {
				assert.Equal(t, 1, len(logs))
			}
			// reset everything
			newNodeRetrieved = false
			nodesAdded = []string{}
			logs = []string{}
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
	FilterPage := func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil }
	endpoint := "https://en.wikipedia.org/wiki/String_cheese"

	// mock out log.Fatalf

	// mute warnings
	originLogWarn := logWarn
	defer func() { logWarn = originLogWarn }()
	logWarn = func(format string, args ...interface{}) {}

	originLogMsg := logMsg
	defer func() { logFatal = originLogMsg }()
	logs := []string{}
	logMsg = func(format string, args ...interface{}) {
		if len(args) > 0 {
			logs = append(logs, fmt.Sprintf(format, args))
		} else {
			logs = append(logs, format)
		}
	}

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
		Crawl("https://en.wikipedia.org/wiki/String_cheese", 2, 1, 0, isValidCrawlLink, addEdges, FilterPage)
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
			1,
			0,
			isValidCrawlLink,
			func(currNode string, neighborNodes []string) ([]string, error) {
				temp := []string{}
				for _, v := range neighborNodes {
					temp = append(temp, "https://en.wikipedia.org"+v)
				}
				nodesAdded = append(nodesAdded, temp...)
				return temp, nil
			},
			func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
		)

		assert.Equal(t, "starting at ["+endpoint+"]", logs[0])
		// only add first recursion nodes, ~30,000 on second recursion
		assert.Equal(t, len(nodesAdded) >= 50, true)
		assert.Equal(t, len(nodesAdded) <= 30000, true)
		assert.Equal(t, []string{}, errors)
	})

	t.Run("adds second level nodes as well", func(t *testing.T) {
		nodesAdded = []string{}
		logs = []string{}
		errors = []string{}
		Crawl(
			endpoint,
			1000,
			1,
			0,
			isValidCrawlLink,
			func(currNode string, neighborNodes []string) ([]string, error) {
				temp := []string{}
				for _, v := range neighborNodes {
					temp = append(temp, "https://en.wikipedia.org"+v)
				}
				nodesAdded = append(nodesAdded, temp...)
				return temp, nil
			},
			func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
		)

		assert.Equal(t, "starting at ["+endpoint+"]", logs[0])
		// only add first recursion nodes, ~30,000 on second recursion
		assert.Equal(t, len(nodesAdded) >= 50, true)
		assert.Equal(t, len(nodesAdded) <= 30000, true)
		assert.Equal(t, []string{}, errors)

	})
	t.Run("on error works correctly", func(t *testing.T) {
		nodesAdded = []string{}
		logs = []string{}
		errors = []string{}
		Crawl(
			endpoint+"/thisisabadendpoint",
			1000,
			1,
			0,
			isValidCrawlLink,
			func(currNode string, neighborNodes []string) ([]string, error) {
				temp := []string{}
				for _, v := range neighborNodes {
					temp = append(temp, "https://skldlfjlskjdflkjsdf.org"+v)
				}
				nodesAdded = append(nodesAdded, temp...)
				return temp, nil
			},
			func(e *colly.HTMLElement) (*colly.HTMLElement, error) { return e, nil },
		)
		assert.Equal(t, 0, len(nodesAdded))
		assert.Equal(t, 1, len(errors))

	})
}
