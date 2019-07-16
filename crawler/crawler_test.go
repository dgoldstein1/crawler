package crawler

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"fmt"
)

func TestCrawl(t *testing.T) {
	isValidCrawlLink := func(url string) bool {
		return strings.HasPrefix(url, "/wiki/") && !strings.Contains(url, ":")
	}
	nodesAdded := []string{}
	addEdges := func(currNode string, neighborNodes []string) ([]string, error) {
		nodesAdded = append(nodesAdded, neighborNodes...)
		return neighborNodes, nil
	}
	connectToDB := func() error { return nil }
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


	t.Run("works with isValidCrawlLink", func(t *testing.T) {
		nodesAdded = []string{}
		// function doing setup of tests
		Crawl("https://en.wikipedia.org/wiki/String_cheese", 2, isValidCrawlLink, connectToDB, addEdges)
		t.Run("only filters on links starting with regex", func(t *testing.T) {
			for _, url := range nodesAdded {
				assert.Equal(t, strings.HasPrefix(url, "/wiki/"), true)
			}
		})
		t.Run("only filters on links starting with regex", func(t *testing.T) {
			for _, url := range nodesAdded {
				if strings.Contains(url, ":") {
					t.Errorf("Did not expect '%s' to contain ':'", url)
				}
			}
		})
	})

	t.Run("adds nodes correctly", func (t *testing.T)  {
		nodesAdded = []string{}
		Crawl(
			endpoint,
			2,
			isValidCrawlLink,
			connectToDB,
			addEdges,
		)

		assert.Equal(t, "starting at [" + endpoint + "]", logs[0])
		// only add first recursion nodes, ~30,000 on second recursion
		assert.Equal(t, len(nodesAdded) > 100, true)
	})
}
