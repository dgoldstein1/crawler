package wikipedia

import (
  "testing"
  "os"
  "fmt"
  "github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T)  {

}

func TestisValidCrawlLink(t *testing.T) {
  t.Run("does not crawl on links with ':'", func(t *testing.T) {
    assert.Equal(t, isValidCrawlLink("/wiki/Category:Spinash"), false)
    assert.Equal(t, isValidCrawlLink("/wiki/Test:"), false)
  })
  t.Run("does not crawl on links not starting with '/wiki/'", func(t *testing.T ){
    assert.Equal(t, isValidCrawlLink("https://wikipedia.org"), false)
    assert.Equal(t, isValidCrawlLink("/wiki"), false)
    assert.Equal(t, isValidCrawlLink("wikipedia/wiki/"), false)
    assert.Equal(t, isValidCrawlLink("/wiki/binary"), true)
  })
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
