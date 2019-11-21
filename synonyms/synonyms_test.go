package synonyms

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var dbEndpoint = "http://localhost:17474"
var twoWayEndpoint = "http://localhost:17475"

func TestIsValidCrawlLink(t *testing.T) {
	t.Run("does not crawl on links with ':'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/synonyms/Category:Spinash"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms/Test:"), false)
	})
	t.Run("does not crawl on links not starting with '/synonyms/'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://synonymspedia.org"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms"), false)
		assert.Equal(t, IsValidCrawlLink("synonymspedia/synonyms/"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms/binary"), true)
	})
	t.Run("doesn't accept 'main_page'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/synonyms/Main_Page"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms/main_Page"), false)
		assert.Equal(t, IsValidCrawlLink("/synonyms/main_page"), false)

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
			URL:              "/synonyms/Maytag_Blue_cheese",
			expectedResponse: "maytag blue cheese",
		},
		Test{
			Name:             "decodes URL in string",
			URL:              "/synonyms/ingeni%c3%b8ren",
			expectedResponse: "ingeniøren",
		},
		Test{
			Name:             "invalid unescape sequence",
			URL:              "/synonyms/^#$%#$G#$(JG#($JG(DFS(J#(JF%23423",
			expectedResponse: "",
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, CleanUrl(test.URL), test.expectedResponse)
		})
	}

}

func TestGetRandomArticle(t *testing.T) {
	errorsLogged := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errorsLogged = append(errorsLogged, fmt.Sprintf(format, args))
		} else {
			errorsLogged = append(errorsLogged, format)
		}
	}

	type Test struct {
		Name             string
		MockedRequest    string
		ExpectedResponse string
		ExpectedError    string
	}

	testTable := []Test{}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// mock out endpoint
			httpmock.Activate()
			httpmock.RegisterResponder("GET", "metawikiEndpoint",
				httpmock.NewStringResponder(200, test.MockedRequest))
			// run test
			a, err := GetRandomArticle()
			assert.Equal(t, test.ExpectedResponse, a)
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), test.ExpectedError))
				assert.Equal(t, 1, len(errorsLogged))
			} else {
				assert.Equal(t, "", test.ExpectedError)
				assert.Equal(t, 0, len(errorsLogged))
			}
			// reset
			httpmock.DeactivateAndReset()
			errorsLogged = []string{}
			timeout = time.Duration(5 * time.Second)
		})
	}

}

func TestAddEdgesIfDoNotExist(t *testing.T) {
	node := "/synonyms/test"
	neighbors := []string{}
	added, _ := AddEdgesIfDoNotExist(node, neighbors)
	assert.Equal(t, added, []string(nil))
}
