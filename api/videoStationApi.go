package api

import (
	"github.com/go-chi/chi"
	param "github.com/oceanicdev/chi-param"
	"net/http"
	"synology-videostation-reindexer/apiModel/response"
)

func (s Server) videoAPIRoutes() {
	s.router.Route("/video", func(r chi.Router) {
		s.addAuthIfNeeded(r)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", s.listLibraries)
			r.Get("/{library}", s.listShares)
			r.Get("/{library}/{share}/reindex", s.reIndex)
		})
	})
}

func (s Server) listShares(w http.ResponseWriter, r *http.Request) {
	library, err := param.String(r, "library")
	if err != nil {
		response.Render(w, r, response.ErrBadRequest(err))
		return
	}
	resp, err := s.videoAPI.ListSharesIn(library)
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.StringArray(resp))
}

func (s *Server) listLibraries(w http.ResponseWriter, r *http.Request) {
	resp, err := s.videoAPI.ListLibraries()
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.StringArray(resp))
}

func (s *Server) reIndex(w http.ResponseWriter, r *http.Request) {
	share, err := param.String(r, "share")
	if err != nil {
		response.Render(w, r, response.ErrBadRequest(err))
		return
	}

	err = s.videoAPI.ReIndexShare(share)
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.String("started reindexing!"))
}
