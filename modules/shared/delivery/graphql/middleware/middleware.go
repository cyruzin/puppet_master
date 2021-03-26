package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cyruzin/puppet_master/domain"
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

type authError struct {
	Error string `json:"error"`
}

func decodeError(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Msg(err.Error())

	e := &authError{Error: err.Error()}

	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.Write([]byte("could not encode the payload"))
		return
	}
}

// AuthMiddleware checks if the request contains Bearer Token on the
// headers and if it is valid.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthentication := r.Header.Get("X-Login")

		if isAuthentication != "" {
			// TODO: Get this value from the config file
			if isAuthentication != "Puppet" {
				decodeError(w, r, errors.New("x-login header do not match"))
				return
			}

			log.Info().Msg("new login, no check is needed")
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			decodeError(w, r, errors.New("authorization header was not provided"))
			return
		}

		// Checking if the header contains Bearer string and if the token exists.
		if !strings.Contains(authHeader, "Bearer") || len(strings.Split(authHeader, "Bearer ")) == 1 {
			decodeError(w, r, errors.New("malformed token"))
			return
		}

		// Capturing the token.
		jwtString := strings.Split(authHeader, "Bearer ")[1]

		token, err := jwt.ParseString(jwtString)
		if err != nil {
			decodeError(w, r, errors.New("failed to parse the token"))
			return
		}

		// Passing permissions through context
		// TODO: Replace the string for the permission struct
		ctx := context.WithValue(r.Context(), domain.ContextKeyID, token.Subject())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
