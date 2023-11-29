package auth

import (
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *Admin)

type AuthMiddleware struct {
	handler AuthenticatedHandler
}

func (am *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	admin := Admin{}
	err := false
	if err {
		// If not, redirect to the login page
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// If they are, call the wrapped handler
	am.handler(w, r, &admin)
}

func Authenticated(handler AuthenticatedHandler) *AuthMiddleware {
	return &AuthMiddleware{handler}
}
