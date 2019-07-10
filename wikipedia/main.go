package wikipedia

import(
  "strings"
  "os"
  "log"
)

// checks environment for required env vars
var logFatalf = log.Fatalf
func parseEnv() {
  requiredEnvs := []string{
    "GRAPH_DB_ENDPOINT",
    "STARTING_ENDPOINT",
    "MAX_CRAWL_DEPTH",
  }
  for _, v := range requiredEnvs{
    if (os.Getenv(v) == "") {
      log.Fatalf("'%s' was not set", v)
    }
  }
}

func isValidCrawlLink(link string) bool {
  return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

func main() {

}
