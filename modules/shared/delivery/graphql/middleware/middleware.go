package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/util"
	"github.com/lestrrat-go/jwx/jwt"
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
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPublic := r.Header.Get("X-Public")

		if isPublic != "" {
			log.Info().Msg("public route, no check is needed")
			ctx := context.WithValue(r.Context(), domain.ContextKeyID, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		isAuthentication := r.Header.Get("X-Login")

		if isAuthentication != "" {
			// TODO: Get this value from the config file
			if isAuthentication != "Puppet" {
				util.DecodeError(w, r, errors.New("x-login header do not match"))
				return
			}

			log.Info().Msg("new login, no check is needed")
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			util.DecodeError(w, r, errors.New("authorization header was not provided"))
			return
		}

		// Checking if the header contains Bearer string and if the token exists.
		if !strings.Contains(authHeader, "Bearer") || len(strings.Split(authHeader, "Bearer ")) == 1 {
			util.DecodeError(w, r, errors.New("malformed token"))
			return
		}

		// Capturing the token.
		jwtString := strings.Split(authHeader, "Bearer ")[1]

		token, err := jwt.ParseString(jwtString)
		if err != nil {
			util.DecodeError(w, r, errors.New("failed to parse the token"))
			return
		}

		// Passing permissions through context
		// TODO: Replace the string for the permission struct
		ctx := context.WithValue(r.Context(), domain.ContextKeyID, token.Subject())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
