package mango

import (
	"flag"
	"fmt"
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

// Prepare all for start
// Separated func for easy testing
func (srv *Server) preStart() {

	// Set default routes
	r := srv.Router

	// doesn't overwrites if user defined same before
	r.HandleFunc("/", srv.runIndex)

	if route := srv.App.URLTemplates["Page"]; route != "" {
		r.HandleFunc(route, srv.runOne)
	}

	if route := srv.App.URLTemplates["Group"]; route != "" {
		r.HandleFunc(route, srv.runGroup)
	}

	// r.HandleFunc("/{slug:[a-z0-9\\-]+}", srv.runOne)
	r.PathPrefix("/{file:.+\\..+}").Handler(http.FileServer(http.Dir(srv.App.PublicPath)))
	r.NotFoundHandler = http.HandlerFunc(srv.run404)

	http.Handle("/", r)

	srv.Templates = template.Must(template.New("").
		Funcs(defaultFuncMap). // fill with defaults
		Funcs(srv.FuncMap).    // user adds/overwrites his own
		ParseGlob(srv.App.BinPath() + "/templates/*.tmpl"))
}

// Start listening to port (default)
func (srv *Server) Start() error {
	srv.preStart()

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
	slug := vars["Slug"]

	page := srv.App.Page(slug)
	if page == nil {
		srv.run404(w, r)
		return
	}

	// Is group
	if page.IsDir() {
		srv.runGroup(w, r)
		return
	}

	// Redirect detected by param
	if redirectURL := page.Get("Redirect"); redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}

	srv.Render(w, page, "one")
}

// Group
func (srv *Server) runGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["Slug"]

	fmt.Println(slug)

	page := srv.App.Page(slug)
	if page == nil || !page.IsDir() {
		srv.run404(w, r)
		return
	}

	srv.Render(w, page, "group")
}

// 404
func (srv *Server) run404(w http.ResponseWriter, r *http.Request) {
	page := srv.App.NewPage("404")
	w.WriteHeader(http.StatusNotFound)
	srv.Render(w, page, "404")
}

// Render only layout
// But give param for page to distinct template
func (srv *Server) Render(w http.ResponseWriter, page *Page, templateID string) {
	page.Set("Template", templateID)
	srv.Templates.ExecuteTemplate(w, "layout", page)
}
