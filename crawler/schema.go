package crawler

// add edge to graph in DB
type addEdgeFunction func(string, string)(bool, error)
// establishes initial connection to DB
type connectToDBFunction func() error
// check if valid url string for crawling
type isValidCrawlLinkFunction func(string) bool
