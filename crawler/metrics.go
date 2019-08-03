package crawler

import (
	"sync/atomic"
)

// constants
var nodesVisited = asyncInt(0)

// resgisters and serves metrics to HTTP
func ServeMetrics() {

}

// increments all counters and gauges
func UpdateMetrics(numberOfNodesAdded int, currDepth int) {
	// increment number of nodes
	nodesVisited.incr(int32(numberOfNodesAdded))
	// increment number of nodes crawled
	// set max depth if greater
}

// increments async int by "n"
func (c *asyncInt) incr(n int32) int32 {
	return atomic.AddInt32((*int32)(c), n)
}

// decrement astnc int
func (c *asyncInt) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
