package synonyms

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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

	testTable := []Test{
		Test{
			Name:          "ENGLISH_WORD_LIST_PATH not set",
			ExpectedError: "ENGLISH_WORD_LIST_PATH was not set",
			Before: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", "")
			},
			After: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", "english.txt")
			},
		},
		Test{
			Name:          "gets random word succesfully",
			ExpectedError: "",
			Before: func() {
				os.Setenv("ENGLISH_WORD_LIST_PATH", "english.txt")
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
				assert.NotEqual(t, w, "")
				assert.Equal(t, err, nil)
			} else {
				assert.Equal(t, w, "")
				assert.NotEqual(t, err, nil)
				assert.Equal(t, err.Error(), test.ExpectedError)
			}
			test.After()
		})
	}

}

func TestAddEdgesIfDoNotExist(t *testing.T) {
	node := "/synonyms/test"
	neighbors := []string{}
	added, _ := AddEdgesIfDoNotExist(node, neighbors)
	assert.Equal(t, added, []string(nil))
}
