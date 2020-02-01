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
	testTable := []struct {
		name              string
		input             string
		expectedToBeValid bool
	}{
		{"positive test", "https://en.wikipedia.org/wiki/Albemarle_County,_Virginia", true},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedToBeValid, IsValidCrawlLink(test.input))
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

	defaultTextDir := "counties.txt"
	testTable := []Test{
		Test{
			Name:          "COUNTIES_LIST not set",
			ExpectedError: "COUNTIES_LIST was not set",
			Before: func() {
				os.Setenv("COUNTIES_LIST", "")
			},
			After: func() {
				os.Setenv("COUNTIES_LIST", defaultTextDir)
			},
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

func TestCleanURL(t *testing.T) {
	type Test struct {
		Name             string
		URL              string
		expectedResponse string
	}

	testTable := []Test{
		Test{
			Name:             "removes prefixes and spaces",
			URL:              "/wiki/Maytag_Blue_cheese",
			expectedResponse: "maytag blue cheese",
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
			url:                    "https://en.wikipedia.org/wiki/Albemarle_County,_Virginia",
			Synonyms:               []string{"Greene County, Virginia"},
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
