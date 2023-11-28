package auth

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

type AdminForm struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required"`
}

type Admin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string
}

type AdminResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type JWTResponse struct {
	Token      string        `json:"token"`
	ExpiresAt  time.Time     `json:"expiresAt"`
	JwtPayload AdminResponse `json:"user"`
}

func createAdminTable(database *sql.DB, logger *lib.Logger) {
	query := `
		CREATE TABLE IF NOT EXISTS admin (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`
	_, err := database.Exec(query)
	if err != nil {
		logger.Error("Error creating admin table" + err.Error())
		return
	}
	logger.Info("Admin table created")
}

func getAdminByUserName(database *sql.DB, username string) (Admin, error) {
	query := `
		SELECT username, email, password
		FROM Admin
		WHERE username= $1
	`
	tx, err := database.Begin()
	if err != nil {
		tx.Rollback()
		return Admin{}, err
	}
	row := database.QueryRow(query, username)
	var admin Admin
	err = row.Scan(&admin.Username, &admin.Email, &admin.Password)
	return admin, err
}

func signUpHandler(database *sql.DB, logger *lib.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var admin AdminForm
		err := json.NewDecoder(r.Body).Decode(&admin)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Invalid form"))
			return
		}
		isValid, message := lib.ValidateForm(admin)

		if !isValid {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(400, message))
			return
		}

		w.WriteHeader(http.StatusOK)
		hashedPassword := lib.HashPassword(admin.Password)
		tx, err := database.Begin()
		if err != nil {
			tx.Rollback()
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(500, "Error creating transaction"+err.Error()))
			return
		}

		checkExistQuery := `
			SELECT username, email
			FROM Admin
			WHERE username = $1 OR email = $2`
		existingAdmin := tx.QueryRow(checkExistQuery, admin.Username, admin.Email)
		var existingUsername string
		var existingEmail string
		err = existingAdmin.Scan(&existingUsername, &existingEmail)

		if err == nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Username or email already exists"))
			return
		}
		query := `
			INSERT INTO admin (username, email, password)
			VALUES ($1, $2, $3)
			RETURNING username, email
		`
		row := tx.QueryRow(query, admin.Username, admin.Email, string(hashedPassword))
		err = row.Scan(&admin.Username, &admin.Email)
		if err != nil {
			tx.Rollback()
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(500, "Error creating admin"+err.Error()))
			return
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(500, "Error creating admin"+err.Error()))
			return
		}
		json.NewEncoder(w).Encode(lib.NewDataResponse(200, "Success", AdminResponse{
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
		admin, err := getAdminByUserName(database, creds.Username)
		isMatch := lib.CheckHashAndPassword(creds.Password, admin.Password)
		if err != nil || !isMatch {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(401, "Invalid credentials"))
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
			JwtPayload: AdminResponse{
				Username: creds.Username,
				Email:    admin.Email,
			},
		})
	}
}

func InitAuthRouter(mux *mux.Router, database *sql.DB, logger *lib.Logger) {
	createAdminTable(database, logger)
	router := mux.PathPrefix("/api/auth").Subrouter()
	router.HandleFunc("/signup", signUpHandler(database, logger)).Methods("POST")
	router.HandleFunc("/signin", signInHandler(database, logger)).Methods("POST")
}
