package crawler

import (
	"testing"
	"github.com/jarcoal/httpmock"
	"net/http"
	"encoding/json"
	"regexp"
	"os"
)

func TestCrawl(t *testing.T) {
	// function doing setup of tests
	t.Run("only filters on links starting with regex", func (t *testing.T)  {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			dbEndpoint := "http://localhost:17474"
			os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
			// connect to DB endpoint
			httpmock.RegisterResponder("GET", dbEndpoint + "/metrics",
				httpmock.NewStringResponder(200, `TEST`))
			// add to DB endpoints
			httpmock.RegisterResponder("GET", dbEndpoint + "/neighbors?node=2",
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, []string{"5","3","7"})
				},
			)
			httpmock.RegisterResponder("POST", dbEndpoint + "/neighbors?node=2",
				func(req *http.Request) (*http.Response, error) {
					body := make(map[string][]string)
					err := json.NewDecoder(req.Body).Decode(&body);
					if err != nil {
						t.Error(err)
					}
					return httpmock.NewJsonResponse(200, body)
				},
			)
		r, _ := regexp.Compile("\\A/wiki/")
		Crawl("https://en.wikipedia.org/wiki/String_cheese", r, 1)
	})
}
