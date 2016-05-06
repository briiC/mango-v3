# Mango [![Go Report Card](http://goreportcard.com/badge/bitbucket.org/briiC/mango-v3)](http://goreportcard.com/report/bitbucket.org/briiC/mango-v3) ![Code coverage](https://img.shields.io/badge/coverage-98.8%-f39f37.svg) [![Docs](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/bitbucket.org/briiC/mango-v3) ![License](https://img.shields.io/badge/license-MIT-blue.svg)

Serves markdown content as webpage.


## Content structure

You can see content structure example in `test-files/`.  
Or find runnable example in folder `example/`.

```
main.go
.mango
content/
    en/
        top-menu/
            1_Home.md
            2_About.md
            ...
        footer-menu/
            ...
    lv/
        ...
public/
    favicon.png
    images/
    css/
    js/
```

## `.mango` - config file
```
Domain: https://example.loc
ContentPath: content/
PublicPath: public/

PageURL: /{Lang}/{Slug}.html
FileURL: /{File}

```

# Examples

Check out `example/README`

## Example (simple)
First install with `go get bitbucket.org/briiC/mango-v3`  
and use in code as `import "bitbucket.org/briiC/mango-v3"`

One-liner if you need basic webpage functionality.

```
#!go
package main

import "bitbucket.org/briiC/mango-v3"

func main() {
    mango.NewServer().Start()
}
```


## Example (advanced)
Before starting webserver add custom stuff if you need advanced configuration.
```
#!go
package main

import "bitbucket.org/briiC/mango-v3"

func main() {
    srv := mango.NewServer()

	// Add some middlewares ("File", "Page")
    // ma.Middlewares["Page"] = mwForPage  // assign one mw
    // ma.Middlewares["File"] = mwForFile  // assign one mw
	srv.Middlewares["Page"] = func(next http.Handler) http.Handler {
		return mwFirst(mwSecond(next))
	}

	// Custom functions for templates
	ma.FuncMap = template.FuncMap{
		"Smile": func() string {
			return ":)"
		},
		"Add": func(a,b int) string {
			return a + b
		},
	}

    // Custom route
	ma.Router.HandleFunc("/{Lang}/search/{sterm}", func(w http.ResponseWriter, r *http.Request) {
		runSearch(ma, w, r) // create your handler
	})

    // Print all pages and info about them
	ma.App.Print()

    // Go!
	log.Println("Start listening on", ":"+ma.Port)
	panic( srv.Start() )
}

func mwFirst(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do stuff
		next.ServeHTTP(w, r)
	})
}
func mwSecond(next http.Handler) http.Handler { ... }
```
