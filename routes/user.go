package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"nojoke/lib"
	"strconv"

	"github.com/gorilla/mux"

	faker "github.com/bxcodec/faker/v3"
	"github.com/gookit/validate"
)

type User struct {
	Id         int    `json:"id"`
	FirstName  string `json:"first_name" validate:"required|minLen:3|maxLen:20"`
	LastName   string `json:"last_name" validate:"required|minLen:3|maxLen:20"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email" validate:"required|email"`
	Age        int    `json:"age" validate:"min:18|max:60"`
	Phone      string `json:"phone"`
	Password   string `json:"password" validate:"minLen:8"`
	Image      string `json:"image"`
}

func (user *User) String() string {
	return user.FirstName + " " + user.LastName
}

func GenerateUsers(limit int) []User {
	userList := []User{}
	for i := 0; i < limit; i++ {
		User := User{}
		User.Id = i
		User.FirstName = faker.FirstName()
		User.LastName = faker.LastName()
		User.MiddleName = faker.FirstName()
		User.Email = faker.Email()
		User.Age = rand.Intn(40) + 20
		User.Password = faker.Password()
		User.Image = faker.URL()
		User.Phone = faker.Phonenumber()
		userList = append(userList, User)
	}
	return userList
}

func validateUserForm(userForm User) (bool, string) {
	v := validate.Struct(userForm)
	if !v.Validate() {
		message := v.Errors.One()
		return false, message
	}
	return true, ""
}

func handleGet(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	total := r.URL.Query().Get("total")

	limitInt, pageInt, totalInt, error := lib.PaginationParams(limit, page, total)

	w.Header().Set("Content-Type", "application/json")

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, "Invalid query params", nil))
		return
	}

	userList := GenerateUsers(totalInt)
	users := lib.PaginateData(userList, limitInt, pageInt, totalInt)

	response := lib.Response{
		Status:  200,
		Message: "OK",
		Data:    users,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	data := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, err.Error(), nil))
		return
	}

	isValid, message := validateUserForm(data)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, message, nil))
		return
	}
	userList := GenerateUsers(100)
	found := false
	for i, user := range userList {
		if user.Id == data.Id || fmt.Sprint(user.Id) == id {
			userList[i] = data
			found = true
			break
		}
	}
	if !found || data.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(lib.NewResponse(404, "User not found", nil))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewResponse(200, "OK", data),
	)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, err.Error(), nil))
		return
	}
	isValid, message := validateUserForm(data)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, message, nil))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewResponse(201, "OK", data),
	)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, "Invalid Id", nil))
		return
	}

	if intId == 0 || intId > 100 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(lib.NewResponse(404, "User not found", nil))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewResponse(200, "DELETED", nil),
	)
}

func handleFindOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewResponse(400, "Invalid Id", nil))
	}
	if intId == 0 || intId > 100 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(lib.NewResponse(404, "User not found", nil))
		return
	}
	userList := GenerateUsers(2)
	user := userList[0]
	user.Id = intId
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lib.NewResponse(200, "OK", user))
}

func initializeMockData(database *sql.DB) {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			phone VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			age INTEGER NOT NULL,
			image VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		); `
	res, err := database.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}

func InitUserRouter(mux *mux.Router, database *sql.DB) {
	router := mux.PathPrefix("/api/users").Subrouter()
	initializeMockData(database)
	router.HandleFunc("", handleGet).Methods("GET")
	router.HandleFunc("", handlePost).Methods("POST")
	router.HandleFunc("/{id}", handleFindOne).Methods("GET")
	router.HandleFunc("/{id}", handlePut).Methods("PUT")
	router.HandleFunc("/{id}", handleDelete).Methods("DELETE")
}
