package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(t *testing.T) {

}

func TestParseEnv(t *testing.T) {
	// mock out log.Fatalf
	origLogFatalf := logFatalf
	defer func() { logFatalf = origLogFatalf }()
	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	requiredEnvs := []string{
		"GRAPH_DB_ENDPOINT",
		"STARTING_ENDPOINT",
		"MAX_APPROX_NODES",
	}

	for _, v := range requiredEnvs {
		os.Setenv(v, "5")
	}
	// positive test
	parseEnv()
	assert.Equal(t, len(errors), 0)

	for _, v := range requiredEnvs {
		t.Run("it validates "+v, func(t *testing.T) {
			errors = []string{}
			os.Unsetenv(v)
			parseEnv()
			assert.Equal(t, len(errors) > 0, true)
			// cleanup
			os.Setenv(v, "5")
		})
	}

	t.Run("fails if MAX_APPROX_NODES is not valid int", func(t *testing.T) {

	})
	t.Run("fails if MAX_APPROX_NODES is not a positive int", func(t *testing.T) {

	})
}
