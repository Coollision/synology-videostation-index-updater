package api

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func (s *Server) initRoutes() {
	r := s.router

	if s.cfg.EnableVideoAPI {
		s.videoAPIRoutes()
	}
	if s.cfg.EnableRadarr {
		s.radarrRoute()
	}
	if s.cfg.EnableSonarr{
		s.sonarrRoute()
	}


	r.Get("/teapot", s.Teapot)
	r.Get("/notimplemented", s.NotImplemented)

}

func (s Server) addAuthIfNeeded(r chi.Router) {
	if s.cfg.Authentication {
		r.Use(middleware.BasicAuth("",
			map[string]string{
				s.cfg.UserName: s.cfg.UserPassword,
			}))
	}
}

// Routes >>>
func (s *Server) Routes() string {
	routes := "synology-videostation-reindexer API\n\n"

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		routes += fmt.Sprintf("%s\t%s\n", method, route)
		return nil
	}

	if err := chi.Walk(s.router, walkFunc); err != nil {
		logrus.Errorf("Failed to walk router: %v", err)
	}

	return routes
}

func (s *Server) registerRoutes() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		if len(route) > 1 && route[len(route)-1] == '/' {
			route = route[:len(route)-1]
		}
		return nil
	}

	if err := chi.Walk(s.router, walkFunc); err != nil {
		logrus.Errorf("Failed to register routes to OAuth: %v", err)
	}
}
