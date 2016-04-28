[![Go Report Card](http://goreportcard.com/badge/bitbucket.org/briiC/mango-v3)](http://goreportcard.com/report/bitbucket.org/briiC/mango-v3) &nbsp;
[![Code coverage](https://img.shields.io/badge/coverage-98.8%-f39f37.svg)](https://img.shields.io/badge/coverage-98.8%-f39f37.svg)  &nbsp;

# Mango
Takes content from markdown files and serves html to browser.


## Content structure
You can see content structure example in `test-files/`.

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


## Example (simple)
One-liner if you need basic webpage functionality.

```
#!go

func main() {
    mango.NewServer().Start()
}
```


## Example (advanced)
Before starting webserver add custom stuff if you need advanced configuration.
```
#!go

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
