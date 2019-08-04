package crawler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	os.Setenv("METRICS_PORT", "8002")
	ServeMetrics()
	res, err := http.Get("http://localhost:8002/metrics")
	require.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	bodyAsString := string(body)
	// assert that has correct values
	assert.True(t, strings.Contains(bodyAsString, "golang_nodes_added"))
	assert.True(t, strings.Contains(bodyAsString, "golang_nodes_visited"))
	assert.True(t, strings.Contains(bodyAsString, "golang_max_depth"))
}

func TestUpdateMetrics(t *testing.T) {
	t.Run("increments nodesVisited", func(t *testing.T) {
		n := totalNodesAdded.get()
		UpdateMetrics(10, 1)
		assert.Equal(t, n+10, totalNodesAdded.get())
	})
	t.Run("set maxDepth if it's greater than current", func(t *testing.T) {
		UpdateMetrics(10, 1)
		assert.Equal(t, int32(1), maxDepth.get())
	})
}
