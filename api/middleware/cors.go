package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// Cors >>>
func Cors() func(http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "Authorization", "Accept-Language"},
		AllowCredentials: true,
	})

	return cors.Handler
}
