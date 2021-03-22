package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LoggerMiddleware logs the details of all requests.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("agent", r.UserAgent()).
			Str("referer", r.Referer()).
			Str("proto", r.Proto).
			Str("remote_address", r.RemoteAddr).
			Dur("latency", time.Since(start)).
			Msg("")

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware checks if the request contains Bearer Token on the
// headers and if it is valid.
// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 1 - Check JWT authenticity (Isolated function)
// 2 - Extract user data (Isolated function)
// 3 - Verify permissions (Isolated function)
//  3.1 - Get all permissions from database or redis (Isolated function)
//  3.2 - Check if the user permissions are valid and return true or false (Isolated function)
// 4 - Return error or forward the request (Here)

// 	})
// }
