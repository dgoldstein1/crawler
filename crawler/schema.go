package crawler

// add edge to graph in DB
// return 'true' if edge already exists
type AddEdgeFunction func(string, []string) ([]string, error)

// establishes initial connection to DB
type ConnectToDBFunction func() error

// check if valid url string for crawling
type IsValidCrawlLinkFunction func(string) bool
