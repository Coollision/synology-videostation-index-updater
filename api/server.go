package api

import (
	"context"
	"fmt"
	"github.com/Coollision/synology-videostation-index-updater/synology/videostation"
	"net/http"
	"time"

	mdw "github.com/Coollision/synology-videostation-index-updater/api/middleware"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type ServerExtension interface {
	AddHandlers(router chi.Router, addAuthIfNeeded func(chi.Router))
}

// Server >>>
type Server struct {
	cfg      *Config
	videoAPI videostation.VideoAPI

	router *chi.Mux
	server *http.Server
}

// NewServer >>>
func NewServer(cfg *Config) *Server {
	server := &Server{
		cfg: cfg,
	}

	server.InitRouter()
	server.InitHTTPServer()

	return server
}

func (s *Server) ImportHandlers(extension ServerExtension) {
	extension.AddHandlers(s.router, s.addAuthIfNeeded)
}

// InitRouter >>>
func (s *Server) InitRouter() {
	s.router = chi.NewRouter()

	s.initMiddleware()
	s.initRoutes()
	s.registerRoutes()
}

func (s *Server) initMiddleware() {
	s.router.Use(
		mdw.Routes("/", s.Routes),
		middleware.Heartbeat("/status"),

		middleware.RequestID,
		middleware.RedirectSlashes,
		middleware.Recoverer,

		mdw.Cors(),

		render.SetContentType(render.ContentTypeJSON),
	)
}

// InitHTTPServer >>>
func (s *Server) InitHTTPServer() {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%v", s.cfg.Port),
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
		Handler:      s.router,
	}
}

// Start >>>
func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

// Stop >>>
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logrus.Errorf("Failed to shutdown server: %v", err)
	}
}
