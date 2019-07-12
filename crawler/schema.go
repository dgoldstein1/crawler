package crawler

// add edge to graph in DB
type AddEdgeFunction func(string, string) (bool, error)

// establishes initial connection to DB
type ConnectToDBFunction func() error

// check if valid url string for crawling
type IsValidCrawlLinkFunction func(string) bool
