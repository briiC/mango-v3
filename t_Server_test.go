package mango

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func Test_Server(t *testing.T) {
	ma := NewServer()
	ma.preStart()

	// Serve to test server
	ts := httptest.NewServer(ma.Router)
	defer ts.Close()

	// Define urls to check
	urls := map[string]map[string]string{
		"/": map[string]string{
			"Code": "200",
			"Body": "</h1>\nindex", // contains
		},
		"/hello": map[string]string{
			"Code": "200",
			"Body": "</h1>\none",
		},
		"/en": map[string]string{
			"Code": "200",
			"Body": "</h1>\ngroup",
		},
		"/en-top-menu": map[string]string{
			"Code": "200",
			"Body": "</h1>\ngroup",
		},
		"/about-cats": map[string]string{ // content from
			"Code": "200",
			"Body": "</h1>\none",
		},
		"/go-to-lv": map[string]string{ // rediret
			"Code": "200",
			"Body": "</h1>\ngroup",
		},
		"/Hello": map[string]string{
			"Code": "404",
			"Body": "404",
		},
	}

	// Check all urls
	for url, m := range urls {
		res, err := http.Get(ts.URL + url)
		if err != nil {
			t.Fatal(url, err)
		}
		body, _ := ioutil.ReadAll(res.Body)

		// HTTP codecheck
		if strconv.Itoa(res.StatusCode) != m["Code"] {
			ma.App.slugPages.Print()
			t.Fatalf("Request[%s] status code should be [%s] not [%d]", url, m["Code"], res.StatusCode)
		}

		// Body check
		if strings.Index(string(body), m["Body"]) == -1 {
			ma.App.slugPages.Print()
			t.Fatalf("Request[%s] body should contain [%s] but found [%s]", url, m["Body"], string(body))
		}
	}

}
