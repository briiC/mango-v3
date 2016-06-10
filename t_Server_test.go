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
	ma := NewServer(3000)
	ma.Middlewares["Page"] = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do smth
			next.ServeHTTP(w, r)
		})
	}
	ma.Middlewares["File"] = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do smth
			next.ServeHTTP(w, r)
		})
	}
	ma.preStart()

	// Serve to test server
	ts := httptest.NewServer(ma.Router)
	defer ts.Close()

	// Define urls to check
	urls := map[string]map[string]string{
		"/": {
			"Code": "200",
			"Body": "</h1>\nindex", // contains
		},
		"/en/hello.html": {
			"Code": "200",
			"Body": "</h1>\none",
		},
		"/en/": {
			"Code": "200",
			"Body": "</h1>\nindex",
		},
		"/en/en-top-menu.html": {
			"Code": "200",
			"Body": "</h1>\ngroup",
		},
		"/en/about-cats.html": { // content from
			"Code": "200",
			"Body": "</h1>\none",
		},
		"/en/-go-to-lv.html": { // redirect
			"Code": "200",
			"Body": "</h1>\nindex",
		},
		"/en/Hello.html": { // gets file 404. case sensitive
			"Code": "404",
			"Body": "404",
		},
		"/en/no-such-file.html": { // gets tmpl 404
			"Code": "404",
			"Body": "404",
		},
	}

	// Check all urls
	for url, m := range urls {
		// fmt.Println(ts.URL + url)
		res, err := http.Get(ts.URL + url)
		if err != nil {
			t.Fatal(url, err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		// fmt.Printf("%s\n", body)

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

	// use for test coverage only
	go ma.Start()

}

//
// func Test_ServerStart(t *testing.T) {
// 	ma := NewServer(3001)
// 	ma.Start()
//
// }
