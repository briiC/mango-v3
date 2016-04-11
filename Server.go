package mango

import (
	"flag"
	"html/template"
	"log"
	"net/http"

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

	return srv
}

// Start listening to port (default)
func (srv *Server) Start() error {

	// Set default routes
	r := srv.Router

	// doesn't overwrites if user defined same before
	r.HandleFunc("/", srv.runIndex)
	r.HandleFunc("/{slug}", srv.runOne)
	r.NotFoundHandler = http.HandlerFunc(srv.run404)
	http.Handle("/", r)

	srv.Templates = template.Must(template.New("").
		Funcs(defaultFuncMap). // fill with defaults
		Funcs(srv.FuncMap).    // user adds/overwrites his own
		ParseGlob(srv.App.BinPath() + "/templates/*.tmpl"))

	// Start listening
	log.Println("Start listening on", ":"+srv.Port)
	return http.ListenAndServe(":"+srv.Port, nil)
}

// Index
func (srv *Server) runIndex(w http.ResponseWriter, r *http.Request) {
	page := srv.App.NewPage("Home")
	srv.Render(w, page, "index")
}

// One
func (srv *Server) runOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	page := srv.App.Page(slug)
	if page == nil {
		srv.run404(w, r)
		return
	}

	srv.Render(w, page, "one")
}

// 404
// TODO: make tru 404
func (srv *Server) run404(w http.ResponseWriter, r *http.Request) {
	page := srv.App.NewPage("404")
	srv.Render(w, page, "404")
}

// Render only layout
// But give param for page to distinct template
func (srv *Server) Render(w http.ResponseWriter, page *Page, templateID string) {
	page.Set("Template", templateID)
	srv.Templates.ExecuteTemplate(w, "layout", page)
}
