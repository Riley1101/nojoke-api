package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nojoke/lib"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var jwtSecret = os.Getenv("JWT_SECRET")

type Admin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type JWTResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Admin     Admin     `json:"user"`
}

func signUpHandler(database *sql.DB, logger *lib.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}

func signInHandler(database *sql.DB, logger *lib.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var creds lib.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expectedPasswd := "admin"
		if creds.Password != expectedPasswd {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Credentials"))
			return
		}
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &lib.Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			logger.Error("" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		json.NewEncoder(w).Encode(JWTResponse{
			Token:     tokenString,
			ExpiresAt: expirationTime,
			Admin: Admin{
				Username: creds.Username,
				Email:    "",
			},
		})
	}
}

func InitAuthRouter(mux *mux.Router, database *sql.DB, logger *lib.Logger) {
	router := mux.PathPrefix("/api/auth").Subrouter()
	router.HandleFunc("/signup", signUpHandler(database, logger)).Methods("POST")
	router.HandleFunc("/signin", signInHandler(database, logger)).Methods("POST")
}
