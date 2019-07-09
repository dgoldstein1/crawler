package crawler

import (
	"testing"
	"strings"
)

func TestCrawl(t *testing.T) {
	t.Run("wikipedia", func(t *testing.T) {
		// function doing setup of tests
		isValidCrawlLink := func(url string) bool {
			return strings.HasPrefix(url, "/wiki/") && !strings.Contains(url, ":")
		}

		nodesAdded := []string{}
		addEdge := func(currNode string, neighborNode string) (bool, error) {
			nodesAdded = append(nodesAdded, neighborNode)
			return false, nil
		}
		connectToDB := func() error {return nil}
		Crawl("https://en.wikipedia.org/wiki/String_cheese", isValidCrawlLink, 2, connectToDB, addEdge)

		t.Run("only filters on links starting with regex", func (t *testing.T)  {
			for _, url := range nodesAdded {
				AssertEqual(t, strings.HasPrefix(url, "/wiki/"), true)
			}
		})
		t.Run("only filters on links starting with regex", func (t *testing.T)  {
			for _, url := range nodesAdded {
				if strings.Contains(url, ":") {
					t.Errorf("Did not expect '%s' to contain ':'", url)
				}
			}
		})
	})
}
