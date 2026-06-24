package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

type WebServer struct {
	Router        chi.Router
	Handlers      []route
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, route{method: method, path: path, handler: handler})
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, r := range s.Handlers {
		s.Router.Method(r.method, r.path, r.handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
