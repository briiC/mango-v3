package mango

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Server - mango server for serving content from mango-structure data
type Server struct {
	Host      string
	Port      string
	App       *Application
	Templates *template.Template
	Router    *mux.Router
	FuncMap   template.FuncMap // user can define its

	/* Middlewares:
	Middlewares["Page"]
	Middlewares["File"]

	To add multiple middlewares on one map key:
	srv.Middlewares["File"] = func(next http.Handler) http.Handler {
		return mwFirst(mwSecond(next))
	}
	*/
	Middlewares map[string]func(next http.Handler) http.Handler
}

// NewServer - create server instance
func NewServer(port int) *Server {

	app, _ := NewApplication()

	srv := &Server{
		Host: "localhost",
		Port: fmt.Sprintf("%d", port),
		App:  app,
	}

	srv.Router = mux.NewRouter()
	srv.Router.StrictSlash(true)

	srv.Middlewares = make(map[string]func(next http.Handler) http.Handler, 0)

	app.Server = srv

	return srv
}

// Prepare all for start
// Separated func for easy testing
func (srv *Server) preStart() http.Handler {

	// Set default routes
	r := srv.Router

	// doesn't overwrites if user defined same before
	r.HandleFunc("/", srv.RunIndex)
	r.HandleFunc("/{Lang:[a-z]{2}}/", srv.RunIndex)

	// Pages (by slug)
	if route := srv.App.URLTemplates["Page"]; route != "" {
		r.HandleFunc(route, srv.RunOne)
	}

	// Files (by file path)
	if route := srv.App.URLTemplates["File"]; route != "" {
		// get prefix to strip
		arr := strings.SplitN(route, "{File", 2)
		prefix := arr[0]

		fs := http.FileServer(http.Dir(srv.App.PublicPath))

		// Middlewares (for files)
		if mw, haveMw := srv.Middlewares["File"]; haveMw {
			fs = mw(fs)
		}

		fs = http.StripPrefix(prefix, fs)
		r.Handle(route, fs)
	}

	// Serve "naked" files. No prefixes, no versions
	// This does nothing if FileURL is: /{File}
	// but mandatory if FileURL is more complex: /static/{File}
	// These lines makes sure we can serve root files: /sitemap.xml
	fs := http.FileServer(http.Dir(srv.App.PublicPath))
	// Middlewares for these files too
	if mw, haveMw := srv.Middlewares["File"]; haveMw {
		fs = mw(fs)
	}
	r.Handle("/{file:.+\\.[a-z]{3,4}}", fs)

	// 404
	r.NotFoundHandler = http.HandlerFunc(srv.Run404)

	// Middlewares (for pages)
	var rh http.Handler
	rh = r
	if mw, haveMw := srv.Middlewares["Page"]; haveMw {
		rh = mw(rh)
	}

	// Try minified templates first
	// If not found use originals
	templatePath := srv.App.BinPath() + "/templates/min"
	if _, err := ioutil.ReadFile(templatePath + "/layout.tmpl"); err != nil {
		templatePath = srv.App.BinPath() + "/templates"
	}
	srv.Templates = template.Must(template.New("").
		Funcs(defaultFuncMap). // fill with defaults
		Funcs(srv.FuncMap).    // user adds/overwrites his own
		ParseGlob(templatePath + "/*.tmpl"))

	return rh
}

// Start listening to port (default)
// Can't be tested because using httptest package (it have his own listener)
func (srv *Server) Start() error {
	rh := srv.preStart()
	http.Handle("/", rh)
	return http.ListenAndServe(":"+srv.Port, nil)
}

// // StartSecure listening to :443 port
// func (srv *Server) StartSecure() error {
// 	rh := srv.preStart()
// 	http.Handle("/", rh)
// 	certPath := srv.App.BinPath() + "/cert.pem"
// 	keyPath := srv.App.BinPath() + "/key.pem"
// 	return http.ListenAndServeTLS(":"+srv.Port, certPath, keyPath, nil)
// }

// RunIndex - handler for first page
func (srv *Server) RunIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lang := vars["Lang"]

	// Default params taken from /content/{lang}/.defaults
	page := srv.App.NewPage(lang, "")
	srv.Render(w, page, "index")
}

// RunOne - handler for specific one (*Page)
func (srv *Server) RunOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["Slug"]

	page := srv.App.Page(slug)
	if page == nil {
		srv.Run404(w, r)
		return
	}

	// Redirect detected by param
	if redirectURL := page.Get("Redirect"); redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}

	templateID := "one"
	if page.IsDir() {
		templateID = "group"
	}
	srv.Render(w, page, templateID)
}

// Run404 - handler 404
func (srv *Server) Run404(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lang := vars["Lang"] // try to get language from url
	page := srv.App.NewPage(lang, "404")
	w.WriteHeader(http.StatusNotFound)

	srv.Render(w, page, "404")
}

// Render only layout
// But give param for page to distinct template
func (srv *Server) Render(w io.Writer, page *Page, templateID string) {
	page.Set("Template", templateID)
	srv.Templates.ExecuteTemplate(w, "layout", page)
}
