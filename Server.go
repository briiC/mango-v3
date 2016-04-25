package mango

import (
	"flag"
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
func NewServer() *Server {

	// Using as string because further code uses it as string
	port := flag.String("port", "3000", "Mango web server port")
	flag.Parse()

	app, _ := NewApplication()

	srv := &Server{
		Host: "localhost",
		Port: *port,
		App:  app,
	}

	srv.Router = mux.NewRouter()
	srv.Router.StrictSlash(true)

	srv.Middlewares = make(map[string]func(next http.Handler) http.Handler, 0)

	return srv
}

// Prepare all for start
// Separated func for easy testing
func (srv *Server) preStart() {

	// Set default routes
	r := srv.Router

	// doesn't overwrites if user defined same before
	r.HandleFunc("/", srv.runIndex)
	r.HandleFunc("/{Lang:[a-z]{2}}", srv.runIndex)

	// Pages (by slug)
	if route := srv.App.URLTemplates["Page"]; route != "" {
		r.HandleFunc(route, srv.runOne)
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

	// 404
	r.NotFoundHandler = http.HandlerFunc(srv.run404)

	// Middlewares (for pages)
	var rh http.Handler
	rh = r
	if mw, haveMw := srv.Middlewares["Page"]; haveMw {
		rh = mw(rh)
	}
	http.Handle("/", rh)

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
}

// Start listening to port (default)
func (srv *Server) Start() error {
	srv.preStart()
	return http.ListenAndServe(":"+srv.Port, nil)
}

// handler: Index
func (srv *Server) runIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lang := vars["Lang"]

	if _, isValid := srv.App.translations[lang]; !isValid {
		// Set default lang if given lang invalid
		lang = srv.App.Pages[0].Get("Slug")
	}

	page := srv.App.NewPage("Home")
	page.Set("Lang", lang)
	srv.Render(w, page, "index")
}

// handler: One (*Page)
func (srv *Server) runOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["Slug"]

	page := srv.App.Page(slug)
	if page == nil {
		srv.run404(w, r)
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

// handler: 404
func (srv *Server) run404(w http.ResponseWriter, r *http.Request) {
	page := srv.App.NewPage("404")
	w.WriteHeader(http.StatusNotFound)
	srv.Render(w, page, "404")
}

// Render only layout
// But give param for page to distinct template
func (srv *Server) Render(w io.Writer, page *Page, templateID string) {
	page.Set("Template", templateID)
	srv.Templates.ExecuteTemplate(w, "layout", page)
}
