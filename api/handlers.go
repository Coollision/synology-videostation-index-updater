package api

import (
	param "github.com/oceanicdev/chi-param"
	"net/http"
	"synology-videostation-reindexer/apiResponses"
)

func (s *Server) Teapot(w http.ResponseWriter, r *http.Request) {
	apiResponses.Render(w, r, apiResponses.ErrTeapot())
}

func (s *Server) NotImplemented(w http.ResponseWriter, r *http.Request) {
	apiResponses.Render(w, r, apiResponses.ErrNotImplemented())
}

func (s *Server) reIndex(w http.ResponseWriter, r *http.Request) {
	library, err := param.String(r, "library")
	if err != nil {
		apiResponses.Render(w, r, apiResponses.ErrBadRequest(err))
		return
	}
	share, err := param.String(r, "share")
	if err != nil {
		apiResponses.Render(w, r, apiResponses.ErrBadRequest(err))
		return
	}

	err = s.videoAPI.ReIndexShare(library, share)
	if err != nil {
		errResp := apiResponses.ErrorToRenderer(err)
		apiResponses.Render(w, r, errResp)
		return
	}
	apiResponses.Render(w, r, apiResponses.String("started reindexing!"))
}

func (s Server) ListShares(w http.ResponseWriter, r *http.Request) {
	library, err := param.String(r, "library")
	if err != nil {
		apiResponses.Render(w, r, apiResponses.ErrBadRequest(err))
		return
	}
	resp, err := s.videoAPI.ListSharesIn(library)
	if err != nil {
		errResp := apiResponses.ErrorToRenderer(err)
		apiResponses.Render(w, r, errResp)
		return
	}
	apiResponses.Render(w, r, apiResponses.StringArray(resp))
}

func (s *Server) ListLibraries(w http.ResponseWriter, r *http.Request) {
	resp, err := s.videoAPI.ListLibraries()
	if err != nil {
		errResp := apiResponses.ErrorToRenderer(err)
		apiResponses.Render(w, r, errResp)
		return
	}
	apiResponses.Render(w, r, apiResponses.StringArray(resp))
}
