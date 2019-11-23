package crawler

import (
	"github.com/gocolly/colly"
)

// add edge to graph in DB
// return 'true' if edge already exists
type AddEdgeFunction func(string, []string) ([]string, error)

// establishes initial connection to DB
type ConnectToDBFunction func() error

// check if valid url string for crawling
type IsValidCrawlLinkFunction func(string) bool

// retrieves new node if current expires
type GetNewNodeFunction func() (string, error)

// number of nodesVisited
type asyncInt int32

// filters page down to more specific element
type FilterPageFunction func(e *colly.HTMLElement) (*colly.HTMLElement, error)
