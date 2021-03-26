package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/util"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			log.Info().Msg("public route, no check is needed")
			ctx := context.WithValue(r.Context(), domain.ContextKeyID, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Checking if the header contains Bearer string and if the token exists.
		if !strings.Contains(authHeader, "Bearer") || len(strings.Split(authHeader, "Bearer ")) == 1 {
			util.DecodeError(w, r, errors.New("malformed token"))
			return
		}

		// Capturing the token.
		jwtString := strings.Split(authHeader, "Bearer ")[1]

		// Parsing the token to verify its authenticity.
		token, err := jwt.ParseString(jwtString, jwt.WithVerify(jwa.HS256, []byte(viper.GetString(`jwt_secret`))))
		if err != nil {
			util.DecodeError(w, r, err)
			return
		}

		// Validating the content.
		if err := jwt.Validate(token); err != nil {
			util.DecodeError(w, r, err)
			return
		}

		// Passing permissions through context
		// TODO: Replace the string for the permission struct
		ctx := context.WithValue(r.Context(), domain.ContextKeyID, token.Subject())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
