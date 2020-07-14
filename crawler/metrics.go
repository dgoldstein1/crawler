package crawler

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
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

	nodesAddedCounter = prometheus.NewCounter(
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
	totalNodesAdded = asyncInt(0)
	maxDepth        = asyncInt(0)
)

// resgisters and serves metrics to HTTP
func ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	// register metrics
	prometheus.MustRegister(nodesVisitedCounter)
	prometheus.MustRegister(nodesAddedCounter)
	prometheus.MustRegister(maxDepthCounter)
	if os.Getenv("METRICS_PORT") == "" {
		os.Setenv("METRICS_PORT", os.Getenv("PORT"))
	}
	// serve http
	go func() {
		logErr("%v", http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("METRICS_PORT")), nil))
	}()
}

// updates prometheus and internal metrics
func UpdateMetrics(numberOfNodesAdded int, currDepth int) {
	// increment number of nodes crawled
	nodesVisitedCounter.Inc()
	// increment number of nodes
	totalNodesAdded.incr(int32(numberOfNodesAdded))
	nodesAddedCounter.Add(float64(numberOfNodesAdded))
	// set max depth if greater

	if int32(currDepth) > maxDepth.get() {
		maxDepth.incr(1)
		maxDepthCounter.Inc()
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
