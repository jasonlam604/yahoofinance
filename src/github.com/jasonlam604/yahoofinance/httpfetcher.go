// Package object httpfetcher implements concurrent utility for
// HTTP retrieval.
//
// Authored by Jason Lam, jasonlam604@gmail.com.
// Code repository hosted on Github at https://github.com/jasonlam604/yahoofinance
// Code is released under MIT License.
package yahoofinance

import (
	"io/ioutil"
	"net/http"
)

// A HttpHandler is callback interface used on success http completion or on error
type HttpHandler interface {
	OnSuccess(httpResponse *HttpResponse)
	OnError(err error)
	OnRequest(url string)
	OnDoneAll()
}

// A HttpResponse represents a http response sent to HttpHandler.OnSuccess when invoked
type HttpResponse struct {
	Url        string
	StatusCode int
	Body       string
	Err        error
}

// Represents client interface
type Connector interface {
	Fetch(HttpHandler, []string) 
}

// httpClient respresents "real" client
type httpClient struct {
}

type Client struct {
	Http Connector
}

// Client Helper createor
func NewHttpFetcher() *Client {
	c := new(Client)
	c.Http = new(httpClient)
	return c
}

// Concurrent HTTP Retriever specific to this package, expects to
// parameters where the first is HttpHandler interface used as callback
// when http processing is done or on error.  Second input parameter is 
// an array of URL strings. No values are returned.
func (r *httpClient) Fetch(httpHandler HttpHandler, urls []string) {

	counter := 0
	ch := make(chan *HttpResponse)
	client := &http.Client{}

	for _, url := range urls {
		go func(url string) {

			req, _ := http.NewRequest( "GET", url, nil )
			
			if resp, err := client.Do(req); nil != err {
				ch <- &HttpResponse{url, http.StatusBadRequest, "", err}
			} else {
				defer resp.Body.Close()
				
				body, _ := ioutil.ReadAll(resp.Body)
				
				ch <- &HttpResponse{url, resp.StatusCode, string(body), err}
			}

		}(url)
	}

	for {
		select {
		case r := <-ch:

			counter++

			httpHandler.OnRequest(r.Url)
			if r.Err != nil {
				httpHandler.OnError(r.Err)
			} else {
				httpHandler.OnSuccess(r)
			}

			if counter == len(urls) {
				httpHandler.OnDoneAll()
				return
			}
		}
	}
}
