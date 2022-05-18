package videostation

import (
	"github.com/Coollision/synology-videostation-index-updater/apiModel/response"
	chi "github.com/go-chi/chi/v5"
	param "github.com/oceanicdev/chi-param"
	"net/http"
)

func (v *videoAPI) AddHandlers(router chi.Router, addAuthIfNeeded func(chi.Router)) {
	router.Route("/video", func(r chi.Router) {
		addAuthIfNeeded(r)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", v.listLibraries)
			r.Get("/{library}", v.listShares)
			r.Get("/{library}/{share}/reindex", v.reIndex)
		})
	})
}

func (v videoAPI) listShares(w http.ResponseWriter, r *http.Request) {
	library, err := param.String(r, "library")
	if err != nil {
		response.Render(w, r, response.ErrBadRequest(err))
		return
	}
	resp, err := v.ListSharesIn(library)
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.StringArray(resp))
}

func (v *videoAPI) listLibraries(w http.ResponseWriter, r *http.Request) {
	resp, err := v.ListLibraries()
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.StringArray(resp))
}

func (v *videoAPI) reIndex(w http.ResponseWriter, r *http.Request) {
	share, err := param.String(r, "share")
	if err != nil {
		response.Render(w, r, response.ErrBadRequest(err))
		return
	}

	err = v.ReIndexShare(share)
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.String("started reindexing!"))
}
