package synonyms

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
		assert.Equal(t, IsValidCrawlLink("/synonym/test"), true)
		assert.Equal(t, IsValidCrawlLink("/synonym/happy"), true)
		assert.Equal(t, IsValidCrawlLink("/synonym/cherry"), true)
	})
	t.Run("does not crawl on links with ':'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/synonym/Category:Spinash"), false)
		assert.Equal(t, IsValidCrawlLink("/synonym/Test:"), false)
	})
	t.Run("does not crawl on links not starting with '/synonym/'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://synonymspedia.org"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms"), false)
		assert.Equal(t, IsValidCrawlLink("synonymspedia/synonym/"), false)
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
			URL:              "/synonym/Maytag_Blue_cheese",
			expectedResponse: "maytag blue cheese",
		},
		Test{
			Name:             "decodes URL in string",
			URL:              "/synonym/ingeni%c3%b8ren",
			expectedResponse: "ingeniøren",
		},
		Test{
			Name:             "invalid unescape sequence",
			URL:              "/synonym/^#$%#$G#$(JG#($JG(DFS(J#(JF%23423",
			expectedResponse: "",
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

	defaultTextDir := "english.txt"
	testTable := []Test{
		Test{
			Name:          "ENGLISH_WORD_LIST_PATH not set",
			ExpectedError: "ENGLISH_WORD_LIST_PATH was not set",
			Before: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", "")
			},
			After: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", defaultTextDir)
			},
		},
		Test{
			Name:          "gets random word succesfully",
			ExpectedError: "",
			Before: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", defaultTextDir)
			},
			After: func() {},
		},
		Test{
			Name:          "no such path",
			ExpectedError: "open this/does/not/exist: no such file or directory",
			Before: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", "this/does/not/exist")
			},
			After: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", defaultTextDir)
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
			DOMLengthMustBeSmaller: 1500,
			url:                    "https://www.synonyms.com/synonym/happy",
			Synonyms:               []string{"felicitous", "glad", "cheerful", "elated"},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create element
			// Request the HTML page.
			res, _ := http.Get(test.url)
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
	node := "/synonym/test"
	neighbors := []string{}
	added, _ := AddEdgesIfDoNotExist(node, neighbors)
	assert.Equal(t, added, []string(nil))
}
