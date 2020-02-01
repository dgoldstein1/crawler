package ar_synonyms

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var dbEndpoint = "http://localhost:17474"
var twoWayEndpoint = "http://localhost:17475"

func TestIsValidCrawlLink(t *testing.T) {
	t.Run("crawls on valid links", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/synonym/ar/test"), true)
		assert.Equal(t, IsValidCrawlLink("/synonym/ar/happy"), true)
		assert.Equal(t, IsValidCrawlLink("/synonym/ar/cherry"), true)
	})
	t.Run("does not crawl on links with ':'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/synonym/ar/Category:Spinash"), false)
		assert.Equal(t, IsValidCrawlLink("/synonym/ar/Test:"), false)
	})
	t.Run("does not crawl on links not starting with '/synonym/ar/'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://synonymspedia.org"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms"), false)
		assert.Equal(t, IsValidCrawlLink("synonymspedia/synonym/ar/"), false)
	})
	t.Run("special use case `https://context.reverso.net/translation/`", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://synonyms.reverso.net/synonym/ar/%D9%86%D9%8A%D8%B3%D8%A7%D9%86"), false)
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

func TestGetRandomNode(t *testing.T) {
	errorsLogged := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errorsLogged = append(errorsLogged, fmt.Sprintf(format, args))
		} else {
			errorsLogged = append(errorsLogged, format)
		}
	}

	type Test struct {
		Name          string
		ExpectedError string
		Before        func()
		After         func()
	}

	defaultTextDir := "arabic.txt"
	testTable := []Test{
		Test{
			Name:          "gets random word succesfully",
			ExpectedError: "",
			Before: func() {
				os.Setenv("ARABIC_WORD_LIST_PATH", defaultTextDir)
			},
			After: func() {},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Before()
			w, err := GetRandomNode()
			// positive tests
			if test.ExpectedError == "" {
				assert.NotEqual(t, "", w)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, "", w)
				assert.NotEqual(t, nil, err)
				assert.Equal(t, test.ExpectedError, err.Error())
			}
			test.After()
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
