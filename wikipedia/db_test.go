package wikipedia

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var dbEndpoint = "http://localhost:17474"
var wikiApiEndpoint = "http://localhost:3000"

func TestIsValidCrawlLink(t *testing.T) {
	t.Run("does not crawl on links with ':'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/wiki/Category:Spinash"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/Test:"), false)
	})
	t.Run("does not crawl on links not starting with '/wiki/'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://wikipedia.org"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki"), false)
		assert.Equal(t, IsValidCrawlLink("wikipedia/wiki/"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/binary"), true)
	})
}

func TestAddToDb(t *testing.T) {
	t.Run("fails when no server found", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
		os.Setenv("WIKI_API_ENDPOINT", wikiApiEndpoint)
		// first test bad response
		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet", []string{"/wiki/Animal"})
		assert.EqualError(t, err, "Get http://localhost:3000?action=parse&format=json&page=Pet&prop=properties: dial tcp 127.0.0.1:3000: connect: connection refused")
		assert.Equal(t, []string{}, newNodes)
	})
	t.Run("returns error when current node doesnt exist (404)", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match

		httpmock.RegisterResponder("GET", wikiApiEndpoint+"?action=parse&format=json&page=Pet_door&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Pet door","pageid":3276454,"properties":[{"name":"wikibase_item","*":"Q943110"}]}}`))
		httpmock.RegisterResponder("GET", wikiApiEndpoint+"?action=parse&format=json&page=Animal&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Animal","pageid":11039790,"properties":[{"name":"wikibase-shortdesc","*":"kingdom of motile multicellular eukaryotic heterotrophic organisms"},{"name":"wikibase_item","*":"Q729"},{"name":"wikibase-badge-Q17437798","*":""}]}}`))
		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(404, map[string]interface{}{
					"code":  404,
					"error": "Node was not found",
				})
			},
		)

		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal"})
		assert.EqualError(t, err, "Node was not found")
		assert.Equal(t, newNodes, []string{})
	})
	t.Run("succesfully adds neighbor nodes", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint+"/edges?node=3276454",
			httpmock.NewStringResponder(200, `[11039790]`))
		httpmock.RegisterResponder("GET", "http://localhost:3000?action=parse&format=json&page=Pet_door&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Pet door","pageid":3276454,"properties":[{"name":"wikibase_item","*":"Q943110"}]}}`))
		httpmock.RegisterResponder("GET", "http://localhost:3000?action=parse&format=json&page=Animal&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Animal","pageid":11039790,"properties":[{"name":"wikibase-shortdesc","*":"kingdom of motile multicellular eukaryotic heterotrophic organisms"},{"name":"wikibase_item","*":"Q729"},{"name":"wikibase-badge-Q17437798","*":""}]}}`))
		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200,  map[string]interface{}{"neighborsAdded" : []int{11039790}})
			},
		)

		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal"})
		assert.Nil(t, err)
		assert.Equal(t, newNodes, []string{"/wiki/Animal"})
	})
	t.Run("only returns new neighbors", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint+"/edges?node=3276454",
			httpmock.NewStringResponder(200, `[11039790]`))
		httpmock.RegisterResponder("GET", "http://localhost:3000?action=parse&format=json&page=Pet_door&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Pet door","pageid":3276454,"properties":[{"name":"wikibase_item","*":"Q943110"}]}}`))
		httpmock.RegisterResponder("GET", "http://localhost:3000?action=parse&format=json&page=Pet_test&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Pet test","pageid":25342,"properties":[{"name":"wikibase_item","*":"Q943110"}]}}`))
		httpmock.RegisterResponder("GET", "http://localhost:3000?action=parse&format=json&page=Animal&prop=properties",
			httpmock.NewStringResponder(200, `{"parse":{"title":"Animal","pageid":11039790,"properties":[{"name":"wikibase-shortdesc","*":"kingdom of motile multicellular eukaryotic heterotrophic organisms"},{"name":"wikibase_item","*":"Q729"},{"name":"wikibase-badge-Q17437798","*":""}]}}`))
		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200,  map[string]interface{}{"neighborsAdded" : []int{11039790}})
			},
		)

		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal", "/wiki/Pet_test"})
		assert.Nil(t, err)
		assert.Equal(t, newNodes, []string{"/wiki/Animal"})
	})
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	t.Run("fails when db not found", func(t *testing.T) {
		err := ConnectToDB()
		assert.EqualError(t, err, "Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused")
	})
	t.Run("succeed when server exists", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint+"/metrics",
			httpmock.NewStringResponder(200, `TEST`))

		err := ConnectToDB()
		assert.Nil(t, err)
	})
}

func TestGetArticleId(t *testing.T) {
	os.Setenv("WIKI_API_ENDPOINT", "https://en.wikipedia.org/w/api.php")
	t.Run("makes request to correct endpoint", func(t *testing.T) {
		id, err := getArticleId("/wiki/Pet")
		assert.Nil(t, err)
		assert.Equal(t, 25079, id)
	})
	t.Run("returns error on bad url", func(t *testing.T) {
		id, err := getArticleId("/wiki/DFSDfet_doorSDFUSFU#UFFISd")
		assert.NotNil(t, err)
		assert.Equal(t, 0, id)
	})

}
