package wikipedia

import (
  "testing"
  "os"
  "fmt"
)

func TestMain(t *testing.T)  {

}

func TestisValidCrawlLink(t *testing.T) {
  t.Run("does not crawl on links with ':'", func(t *testing.T) {
    AssertEqual(t, isValidCrawlLink("/wiki/Category:Spinash"), false)
    AssertEqual(t, isValidCrawlLink("/wiki/Test:"), false)
  })
  t.Run("does not crawl on links not starting with '/wiki/'", func(t *testing.T ){
    AssertEqual(t, isValidCrawlLink("https://wikipedia.org"), false)
    AssertEqual(t, isValidCrawlLink("/wiki"), false)
    AssertEqual(t, isValidCrawlLink("wikipedia/wiki/"), false)
    AssertEqual(t, isValidCrawlLink("/wiki/binary"), true)
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
    AssertEqual(t, len(errors), 0)

    for _, v := range requiredEnvs {
      t.Run("it validates " + v, func (t *testing.T)  {
        errors = []string{}
        os.Unsetenv(v)
        parseEnv()
        AssertEqual(t, len(errors), 1)
        AssertEqual(t, errors[0], "'" + v + "' was not set")
        // cleanup
        os.Setenv(v, "TEST")
      })
    }
}
