package crawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
		Crawl("https://en.wikipedia.org/wiki/String_cheese", 2, isValidCrawlLink, connectToDB, addEdges)
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
			connectToDB,
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
