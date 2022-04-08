package sonarr

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"synology-videostation-reindexer/apiModel/request"
	"synology-videostation-reindexer/apiModel/response"
	"synology-videostation-reindexer/sonarr_radarrhooks"
	"synology-videostation-reindexer/synology/videostation"
)

var log = logrus.WithField("hook", "sonarr")

func NewHook(cfg sonarr_radarrhooks.HooksConfig, api videostation.VideoAPI) *hook {
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config for sonarr %s", err)
	}
	h := &hook{Cfg: cfg, VideoAPI: api}
	h.setReindexFunction()
	return h
}

type hook sonarr_radarrhooks.Hook

func (h *hook) AddHandlers(router chi.Router, addAuthIfNeeded func(chi.Router)) {
	if h.Cfg.Enabled {
		router.Route("/sonarr", func(r chi.Router) {
			addAuthIfNeeded(r)
			r.Post("/", h.sonarr)
		})
	}
}

func (h *hook) sonarr(w http.ResponseWriter, r *http.Request) {
	data := request.Sonarr{}
	err := render.Bind(r, &data)
	if err != nil {
		log.Errorf("got a bad request %e", err)
		response.Render(w, r, response.ErrBadRequest(err))
		return
	}
	log.Debug(data)
	if data.EventType == "Test" {
		log.Infof("got a testRequest")
		response.Render(w, r, response.String("test was successful"))
		return
	}

	if err := h.Reindex(); err != nil {
		log.Errorf("failed to start reindexing: %e", err)
		response.Render(w, r, response.ErrInternalServer(err))
		return
	}
	log.Info("series started reindexing")
	response.Render(w, r, response.String("series reindexing started"))
}

func (h *hook) setReindexFunction() {
	if h.Cfg.Library != "" {
		h.Reindex = func() error {
			shares, err := h.VideoAPI.ListSharesIn(h.Cfg.Library)
			if err != nil {
				return err
			}
			for _, share := range shares {
				err := h.VideoAPI.ReIndexShare(share)
				if err != nil {
					return err
				}
			}
			return nil
		}
	} else {
		h.Reindex = func() error {
			err := h.VideoAPI.ReIndexShare(h.Cfg.Share)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
