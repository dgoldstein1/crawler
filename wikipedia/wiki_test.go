package wikipedia

import (
	"fmt"
	"github.com/gocolly/colly"
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
			URL:              "/wiki/Maytag_Blue_cheese",
			expectedResponse: "maytag blue cheese",
		},
		Test{
			Name:             "decodes URL in string",
			URL:              "/wiki/ingeni%c3%b8ren",
			expectedResponse: "ingeniøren",
		},
		Test{
			Name:             "invalid unescape sequence",
			URL:              "/wiki/^#$%#$G#$(JG#($JG(DFS(J#(JF%23423",
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
		Name             string
		MockedRequest    string
		ExpectedResponse string
		ExpectedError    string
	}

	testTable := []Test{
		Test{
			Name:             "succesful",
			MockedRequest:    `{"batchcomplete":"","continue":{"grncontinue":"0.369259750651|0.369260921533|12247122|0","continue":"grncontinue||"},"query":{"pages":{"9820486":{"pageid":9820486,"ns":0,"title":"Oregon Bicycle Racing Association","extract":"The Oregon Bicycle Racing Association is a bicycle racing organization based in the U.S. state of Oregon."}}},"limits":{"extracts":20}}`,
			ExpectedResponse: "https://en.wikipedia.org/wiki/Oregon Bicycle Racing Association",
			ExpectedError:    "",
		},
	}

	tempEndpointVar := ""
	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// mock out endpoint
			if test.Name == "ENDPOINT_NOT_FOUND" {
				tempEndpointVar = metawikiEndpoint
				metawikiEndpoint = "http://BAD_ENDPOINT"
				timeout = time.Duration(1 * time.Second)
			} else {
				httpmock.Activate()
				httpmock.RegisterResponder("GET", metawikiEndpoint,
					httpmock.NewStringResponder(200, test.MockedRequest))
			}
			// run test
			a, err := GetRandomNode()
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
			metawikiEndpoint = tempEndpointVar
		})
	}

}

func TestFilterPage(t *testing.T) {
	type Test struct {
		Name          string
		ExpectedError string
		ExpectedText  string
		el            colly.HTMLElement
	}

	testTable := []Test{
		Test{
			Name:          "positive test",
			ExpectedError: "",
			ExpectedText:  "test",
			el: colly.HTMLElement{
				Text: "test",
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// run tests
			e, err := FilterPage(&test.el)
			if test.ExpectedError == "" {
				assert.Equal(t, nil, err)
			} else {
				assert.NotEqual(t, nil, err)
			}
			assert.Equal(t, test.ExpectedText, e.Text)
		})
	}

}

func TestAddEdgesIfDoNotExist(t *testing.T) {
	node := "/wiki/test"
	neighbors := []string{}
	added, _ := AddEdgesIfDoNotExist(node, neighbors)
	assert.Equal(t, added, []string(nil))
}
