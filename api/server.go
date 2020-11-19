package api

import (
	"context"
	"fmt"
	"net/http"
	"synology-videostation-reindexer/synology"
	"time"

	mdw "synology-videostation-reindexer/api/middleware"

	"github.com/Remitly/chi-prometheus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Server >>>
type Server struct {
	cfg  *Config
	syno *synology.SynoStuff

	router *chi.Mux
	server *http.Server
}

// NewServer >>>
func NewServer(cfg *Config, syno *synology.SynoStuff) *Server {
	server := &Server{
		cfg:  cfg,
		syno: syno,
	}

	server.InitRouter()
	server.InitHTTPServer()

	return server
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

		chiprometheus.NewMiddleware("synology-videostation-reindexer"),

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
		WriteTimeout: time.Minute * 5, // this is high because a request can take long to write
		ReadTimeout:  time.Minute * 2,
		IdleTimeout:  time.Minute * 2,
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
