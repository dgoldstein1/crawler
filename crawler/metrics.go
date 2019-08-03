package crawler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync/atomic"
)

// constants
var (
	nodesVisitedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "nodes_visited",
			Help:      "Number of nodes scraped and succesfully added to the graph",
		})

	nodesAddedCoutner = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "nodes_added",
			Help:      "Number of nodes succesfully visited",
		})

	maxDepthCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "max_depth",
			Help:      "Max depth in the tree visited nodes",
		})
	nodesVisited = asyncInt(0)
	maxDepth     = asyncInt(0)
)

// resgisters and serves metrics to HTTP
func ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	// register metrics
	prometheus.MustRegister(nodesVisitedCounter)
	prometheus.MustRegister(nodesAddedCoutner)
	prometheus.MustRegister(maxDepthCounter)
	// serve http
	go func() {
		logErr("%v", http.ListenAndServe(":8080", nil))
	}()
}

// updates prometheus and internal metrics
func UpdateMetrics(numberOfNodesAdded int, currDepth int) {
	// increment number of nodes crawled
	nodesVisitedCounter.Inc()
	// increment number of nodes
	nodesVisited.incr(int32(numberOfNodesAdded))
	nodesVisitedCounter.Add(float64(numberOfNodesAdded))
	// set max depth if greater
	if int32(currDepth) > maxDepth.get() {
		maxDepth.incr(int32(currDepth) - maxDepth.get())
		maxDepthCounter.Add(float64(int32(currDepth) - maxDepth.get()))
	}
}

// increments async int by "n"
func (c *asyncInt) incr(n int32) int32 {
	return atomic.AddInt32((*int32)(c), n)
}

// decrement astnc int
func (c *asyncInt) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
