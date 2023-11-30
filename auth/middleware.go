package auth

import (
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *Admin)

type AuthMiddleware struct {
	handler AuthenticatedHandler
}

func (am *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get jwt
	// check if jwt is valid
	// if valid, get admin
	// if not, redirect to guest page
	jwt := r.Header.Get("Authorization")
	admin := Admin{}
	if jwt == "" {
		am.handler(w, r, nil)
		return
	}
	am.handler(w, r, &admin)
}

func Authenticated(handler AuthenticatedHandler) *AuthMiddleware {
	return &AuthMiddleware{handler}
}
