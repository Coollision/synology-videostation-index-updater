package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
)

// Routes >>>
func Routes(endpoint string, routeProvider func() string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		var (
			init   sync.Once
			routes string
		)
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.EqualFold(r.URL.Path, endpoint) {
				init.Do(func() {
					routes = routeProvider()
				})
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(routes))
				if err != nil {
					logrus.Errorf("failed to write response: %v", err)
				}
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
