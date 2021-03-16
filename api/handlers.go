package api

import (
	"net/http"
	"synology-videostation-reindexer/apiModel/response"
)

func (s *Server) Teapot(w http.ResponseWriter, r *http.Request) {
	response.Render(w, r, response.ErrTeapot())
}

func (s *Server) NotImplemented(w http.ResponseWriter, r *http.Request) {
	response.Render(w, r, response.ErrNotImplemented())
}


