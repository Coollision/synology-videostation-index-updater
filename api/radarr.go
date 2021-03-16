package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"synology-videostation-reindexer/apiModel/request"
	"synology-videostation-reindexer/apiModel/response"
)

func (s Server) radarrRoute() {
	s.router.Route("/radarr", func(r chi.Router) {
		s.addAuthIfNeeded(r)
		r.Post("/", s.radarr)
	})
}

func (s *Server) radarr(w http.ResponseWriter, r *http.Request) {
	log:= logrus.WithField("handlers","radarr")
	data:= request.Radarr{}
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
	log.Info("videos started reindexing")
	response.Render(w, r, response.String("videos reindexing started"))
}
