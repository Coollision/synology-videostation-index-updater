package sonarr

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"synology-videostation-reindexer/apiModel/request"
	"synology-videostation-reindexer/apiModel/response"
)

type Hooks struct {}

func (h Hooks) AddHandlers(router chi.Router, addAuthIfNeeded func(chi.Router)) {
	router.Route("/sonarr", func(r chi.Router) {
		addAuthIfNeeded(r)
		r.Post("/", sonarr)
	})
}

func sonarr(w http.ResponseWriter, r *http.Request) {
	log:= logrus.WithField("handlers","sonarr")
	data:= request.Sonarr{}
	err := render.Bind(r, &data)
	if err != nil{
		log.Errorf("got a bad request %e", err)
		response.Render(w, r, response.ErrBadRequest(err))
		return
	}
	if data.EventType == "Test"{
		log.Infof("got a testRequest")
		response.Render(w,r,response.String("test was successful"))
		return
	}
	log.Info("series started reindexing")
	response.Render(w, r, response.String("series reindexing started"))
}
