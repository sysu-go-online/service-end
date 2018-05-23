package middleware

import (
	"net/http"
)

type AuthToken struct{}

func (a AuthToken) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

}
