package crawler

import (
	"testing"
	"regexp"
	"fmt"
)

func TestCrawl(t *testing.T) {
	// function doing setup of tests
	t.Run("only filters on links starting with regex", func (t *testing.T)  {
		r, _ := regexp.Compile("\\A/wiki/")
		nodesAdded := []string{}
		addEdge := func(currNode string, neighborNode string) (bool, error) {
			nodesAdded = append(nodesAdded, neighborNode)
			return false, nil
		}
		connectToDB := func() error {return nil}
		Crawl("https://en.wikipedia.org/wiki/String_cheese", r, 2, connectToDB, addEdge)
		fmt.Println(nodesAdded)
	})
}
