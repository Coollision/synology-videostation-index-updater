package radarr

import (
	"github.com/Coollision/synology-videostation-index-updater/apiModel/request"
	"github.com/Coollision/synology-videostation-index-updater/apiModel/response"
	"github.com/Coollision/synology-videostation-index-updater/sonarr_radarrhooks"
	"github.com/Coollision/synology-videostation-index-updater/synology/videostation"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.WithField("hook", "radarr")

func NewHook(cfg sonarr_radarrhooks.HooksConfig, api videostation.VideoAPI) *hook {
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config for radarr %s", err)
	}
	h := &hook{Cfg: cfg, VideoAPI: api}
	h.setReindexFunction()
	return h
}

type hook sonarr_radarrhooks.Hook

func (h *hook) AddHandlers(router chi.Router, addAuthIfNeeded func(chi.Router)) {
	if h.Cfg.Enabled {
		router.Route("/radarr", func(r chi.Router) {
			addAuthIfNeeded(r)
			r.Post("/", h.radarr)
		})
	}
}

func (h *hook) radarr(w http.ResponseWriter, r *http.Request) {
	data := request.Radarr{}
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
	log.Info("videos started reindexing")
	response.Render(w, r, response.String("videos reindexing started"))
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
