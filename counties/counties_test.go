package ar_synonyms

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var dbEndpoint = "http://localhost:17474"
var twoWayEndpoint = "http://localhost:17475"

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
	t.Run("doesn't accept 'main_page'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/wiki/Main_Page"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/main_Page"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/main_page"), false)

	})
}

func TestCleanURL(t *testing.T) {
	type Test struct {
		Name             string
		URL              string
		expectedResponse string
	}

	testTable := []Test{
		Test{
			Name:             "removes prefixes and spaces",
			URL:              "/synonym/ar/%D8%AD%D9%8A%D9%86",
			expectedResponse: "حين",
		},
		Test{
			Name:             "decodes URL in string",
			URL:              "/synonym/ar/ingeni%c3%b8ren",
			expectedResponse: "ingeniøren",
		},
		Test{
			Name:             "invalid unescape sequence",
			URL:              "/synonym/ar/^#$%#$G#$(JG#($JG(DFS(J#(JF%23423",
			expectedResponse: "",
		},
		Test{
			Name:             "removes 'https' with base endpoint as well",
			URL:              "https://synonyms.reverso.net/synonym/ar/موسم",
			expectedResponse: "موسم",
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, CleanUrl(test.URL), test.expectedResponse)
		})
	}

}

func TestFilterPage(t *testing.T) {
	type Test struct {
		Name                   string
		ExpectedError          string
		DOMLengthMustBeGreater int
		DOMLengthMustBeSmaller int
		Synonyms               []string
		url                    string
	}

	testTable := []Test{
		Test{
			Name:                   "positive test",
			ExpectedError:          "",
			DOMLengthMustBeGreater: 0,
			DOMLengthMustBeSmaller: 40000,
			url:                    "https://synonyms.reverso.net/synonym/ar/%D8%AF%D9%88%D8%B1",
			Synonyms:               []string{"مرحلة"}, // []string{"مرحلة", "فترة", "وقت", "عصر"},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create element
			// Request the HTML page.
			client := &http.Client{}
			req, err := http.NewRequest("GET", test.url, nil)
			// need
			req.Header.Add("User-Agent", `Dgoldstein1/crawler`)
			res, _ := client.Do(req)
			defer res.Body.Close()
			// Load the HTML document
			doc, _ := goquery.NewDocumentFromReader(res.Body)
			el := colly.HTMLElement{
				DOM: doc.Selection,
			}
			// run tests
			e, err := FilterPage(&el)
			if test.ExpectedError == "" {
				assert.Equal(t, nil, err)
			} else {
				assert.NotEqual(t, nil, err)
			}
			assert.Less(t, test.DOMLengthMustBeGreater, len(e.DOM.Text()))
			assert.Greater(t, test.DOMLengthMustBeSmaller, len(e.DOM.Text()))
			// make sure there are href links
			for _, w := range test.Synonyms {
				assert.Contains(t, e.DOM.Find("a[href]").Text(), w)
			}
		})
	}
}

func TestAddEdgesIfDoNotExist(t *testing.T) {
	node := "/synonym/ar/test"
	neighbors := []string{}
	added, _ := AddEdgesIfDoNotExist(node, neighbors)
	assert.Equal(t, added, []string(nil))
}
