// httpfetcher unit test.
//
// Authored by Jason Lam, jasonlam604@gmail.com.
// Code repository hosted on Github at https://github.com/jasonlam604/yahoofinance
// Code is released under MIT License.
package yahoofinance

import (
	"fmt"
	"net/http"
	"testing"
)

// Mock Client
type MockHttpClient struct {
}

// Mock Client Creator
func getTestClient() *Client {
	c := new(Client)
	c.Http = new(MockHttpClient)
	return c
}

// Observer for Http Events/callbacks
type ClientObserver struct {
}

func (c ClientObserver) OnSuccess(httpResponse *HttpResponse) {
	flagtests["success"] = 1
	responsetest = httpResponse
}

func (c ClientObserver) OnError(err error) {
	flagtests["error"] = 1
}

func (c ClientObserver) OnRequest(url string) {
	flagtests["request"] += 1
}

func (c ClientObserver) OnDoneAll() {
	flagtests["alldone"] += 1
}

// Test Flags and test data
var (
	urlstests = []string{
		"http://www.fake-site-success.com",
		"http://www.fake-site-error.com",
		"http://www.fake-site-sucess-trigger-request.com",
		"http://www.fake-site-3.com",
	}

	responsetest = &HttpResponse{}

	flagtests = map[string]int{
		"success": 0,
		"error":   0,
		"request": 0,
		"alldone": 0,
	}
)

// Mockup Fetch
func (mwc *MockHttpClient) Fetch(httpHandler HttpHandler, urls []string) {

	for _, url := range urls {

		httpHandler.OnRequest(url)

		if url == "http://www.fake-site-success.com" {
			r := &HttpResponse{"http://www.fake-site-success.com", http.StatusOK, "fake body data", nil}
			httpHandler.OnSuccess(r)
		}

		if url == "http://www.fake-site-error.com" {
			httpHandler.OnError(nil)
		}
	}

	httpHandler.OnDoneAll()

}

// Unit Tester
func TestFetchEvents(t *testing.T) {

	clientObserver := ClientObserver{}

	c := getTestClient()
	c.Http.Fetch(clientObserver, urlstests)

	fmt.Println(flagtests["request"])

	if flagtests["error"] != 1 {
		t.Errorf("Errors actual %d, expected %d", flagtests["error"], 1)
	}

	if responsetest.Url != "http://www.fake-site-success.com" {
		t.Errorf("Response actual url %s, expected %s", responsetest.Url, "http://www.fake-site-success.com")
	}

	if flagtests["request"] != len(urlstests) {
		t.Errorf("Requests actual %d, expected %d", flagtests["request"], len(urlstests))
	}

	if flagtests["alldone"] != 1 {
		t.Errorf("Requests actual %d, expected %d", flagtests["alldone"], 1)
	}
}
