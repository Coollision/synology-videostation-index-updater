package response

import (
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Render stub to end error handling madness
func Render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		logrus.Errorf("failed to render response: %v", err)
	}
}
