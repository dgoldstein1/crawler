package crawler

import (
	"testing"
	// "regexp"
	"reflect"
	"os"
	"github.com/jarcoal/httpmock"
)


// TODO: implement after connectToDB to addToDB
// func TestCrawl(t *testing.T) {
// 	r, _ := regexp.Compile("\\A/wiki/")
//
// 	Crawl("https://en.wikipedia.org/wiki/String_cheese", r, 2)
// }

func TestAddToDb(t *testing.T) {}

func TestDoStuffWithTestServer(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	// mock out http endpoint
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
		// Exact URL match
	httpmock.RegisterResponder("GET", dbEndpoint + "/metrics",
		httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
	// Use Client & URL from our local test server
	err := connectToDB()
	AssertErrorEqual(t, err, nil)
}

// adopted taken from https://gist.github.com/samalba/6059502
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received '%v' (type %v), expected '%v' (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func AssertErrorEqual(t *testing.T, a error, b error) {
	if (a == nil || b == nil) {
		AssertEqual(t, a, b)
		return
	}
	if (a.Error() == b.Error()) {
		return
	}
	t.Errorf("Received '%v' (type %v), expected '%v' (type %v)", a.Error(), reflect.TypeOf(a), b.Error(), reflect.TypeOf(b))
}
