package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nojoke/lib"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/validate"
	"github.com/gorilla/mux"
)

var jwtSecret = os.Getenv("JWT_SECRET")

type AdminForm struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required"`
}

type Admin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type JWTResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Admin     Admin     `json:"user"`
}

func validateAdminForm(admin AdminForm) (bool, string) {
	v := validate.Struct(admin)
	if !v.Validate() {
		message := ""
		message += v.Errors.One()
		return false, message
	}
	return true, ""
}

func signUpHandler(database *sql.DB, logger *lib.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var admin AdminForm
		err := json.NewDecoder(r.Body).Decode(&admin)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(lib.NewResponse(400, "Invalid form", nil))
			return
		}
		isValid, message := validateAdminForm(admin)

		if !isValid {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(lib.NewResponse(400, message, nil))
			return
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(lib.NewResponse(200, "Success", Admin{
			Username: admin.Username,
			Email:    admin.Email,
		}))
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
				Email:    "admin@gmail.com",
			},
		})
	}
}

func InitAuthRouter(mux *mux.Router, database *sql.DB, logger *lib.Logger) {
	router := mux.PathPrefix("/api/auth").Subrouter()
	router.HandleFunc("/signup", signUpHandler(database, logger)).Methods("POST")
	router.HandleFunc("/signin", signInHandler(database, logger)).Methods("POST")
}
