package main

import (
  "testing"
  "os"
  "fmt"
  "github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T)  {

}

func TestParseEnv(t *testing.T) {
    // mock out log.Fatalf
    origLogFatalf := logFatalf
    defer func() { logFatalf = origLogFatalf } ()
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
      "MAX_CRAWL_DEPTH",
    }

    for _, v := range requiredEnvs {
      os.Setenv(v, "TEST")
    }
    // positive test
    parseEnv()
    assert.Equal(t, len(errors), 0)

    for _, v := range requiredEnvs {
      t.Run("it validates " + v, func (t *testing.T)  {
        errors = []string{}
        os.Unsetenv(v)
        parseEnv()
        assert.Equal(t, len(errors), 1)
        // cleanup
        os.Setenv(v, "TEST")
      })
    }
}
