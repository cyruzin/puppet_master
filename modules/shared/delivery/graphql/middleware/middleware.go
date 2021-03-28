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
		token, err := jwt.ParseString(jwtString, jwt.WithVerify(jwa.HS256, []byte(viper.GetString(`jwt.secret`))))
		if err != nil {
			util.DecodeError(w, r, err)
			return
		}

		// Validating the content.
		if err := jwt.Validate(token); err != nil {
			util.DecodeError(w, r, errors.New("invalid token"))
			return
		}

		userInfo := token.PrivateClaims()["user"].(map[string]interface{})

		userID := int64(userInfo["user_id"].(float64))
		roles := []string{}
		userPermissions := []string{}
		rolePermissions := []string{}
		roleIDs := []int64{}

		for _, role := range userInfo["roles"].([]interface{}) {
			roles = append(roles, role.(string))
		}

		for _, roleID := range userInfo["role_ids"].([]interface{}) {
			roleIDs = append(roleIDs, int64(roleID.(float64)))
		}

		for _, permission := range userInfo["role_permissions"].([]interface{}) {
			rolePermissions = append(rolePermissions, permission.(string))
		}

		for _, permission := range userInfo["user_permissions"].([]interface{}) {
			userPermissions = append(userPermissions, permission.(string))
		}

		auth := &domain.Auth{
			UserID:          userID,
			RoleIDs:         roleIDs,
			Name:            userInfo["name"].(string),
			Email:           userInfo["email"].(string),
			Roles:           roles,
			RolePermissions: userPermissions,
			UserPermissions: rolePermissions,
		}

		// Passing permissions through context
		ctx := context.WithValue(r.Context(), domain.ContextKeyID, auth)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
