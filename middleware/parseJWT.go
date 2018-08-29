package middleware

import (
	"net/http"
)

// ParseJWT parse token from header
// and add add X-Username to the header
type ParseJWT struct{}

func (a ParseJWT) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO:
}
