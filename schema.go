package crawler

// add edge to graph in DB
type addEdgeFunction func(string, string)(bool, error)
// establishes initial connection to DB
type connectToDBFunction func() error
