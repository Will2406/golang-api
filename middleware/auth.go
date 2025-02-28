package middleware

import (
	"golang-api/models"
	"golang-api/server"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"loign",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			CheckAuthorizationHeader(s, w, r)
			next.ServeHTTP(w, r)
		})
	}
}

func CheckAuthorizationHeader(s server.Server, w http.ResponseWriter, r *http.Request) *jwt.Token {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	return token
}
