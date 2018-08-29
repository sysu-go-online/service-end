package middleware

import (
	"net/http"
)

// Auth auth user token
type Auth struct{}

func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO:
}
