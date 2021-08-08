package testing

import (
	"github.com/go-chi/chi"
	"net/http"
	"synology-videostation-reindexer/apiModel/response"
)

func (v *testingAPI) AddHandlers(router chi.Router, addAuthIfNeeded func(chi.Router)) {
	router.Route("/test", func(r chi.Router) {
			r.Get("/", v.testapi)
		})
}



func (v *testingAPI) testapi(w http.ResponseWriter, r *http.Request) {
	resp, err := v.test()
	if err != nil {
		errResp := response.ErrorToRenderer(err)
		response.Render(w, r, errResp)
		return
	}
	response.Render(w, r, response.String(resp))
}
