package crawler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsyncInt(t *testing.T) {
	t.Run("able to increment succesfully", func(t *testing.T) {
		nodesVisited := asyncInt(0)
		nodesVisited.incr(253)
		assert.Equal(t, int32(nodesVisited), int32(253))
	})
	t.Run("able to get succesfully", func(t *testing.T) {
		nodesVisited := asyncInt(25342)
		assert.Equal(t, nodesVisited.get(), int32(25342))
	})
}

func TestServeServiceMetrics(t *testing.T) {

}

func TestUpdateMetrics(t *testing.T) {
	t.Run("increments nodesVisited", func(t *testing.T) {
		n := nodesVisited.get()
		UpdateMetrics(10, 1)
		assert.Equal(t, n+10, nodesVisited.get())
	})
	t.Run("set maxDepth if it's greater than current", func(t *testing.T) {
		UpdateMetrics(10, 100)
		assert.Equal(t, int32(100), maxDepth.get())
	})
}
