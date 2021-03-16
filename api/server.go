package api

import (
	"context"
	"fmt"
	"net/http"
	"synology-videostation-reindexer/synology/videostation"
	"time"

	mdw "synology-videostation-reindexer/api/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
func NewServer(cfg *Config, syno videostation.VideoAPI) *Server {
	server := &Server{
		cfg:      cfg,
		videoAPI: syno,
	}

	server.InitRouter()
	server.InitHTTPServer()

	return server
}


func (s *Server)ImportHandlers(extension ServerExtension){
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
		logrus.Errorf("Failed to start server: %v", err)
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




