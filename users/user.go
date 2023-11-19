package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"nojoke/lib"
	"strconv"

	faker "github.com/bxcodec/faker/v3"
	"github.com/gookit/validate"
	"github.com/gorilla/mux"
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
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Invalid query params"))
		return
	}

	userList := GenerateUsers(totalInt)
	users := lib.PaginateData(userList, limitInt, pageInt, totalInt)

	response := lib.DataResponse{
		Status:  200,
		Message: "OK",
		Data:    users,
		Pagination: lib.Pagination{
			Total: totalInt,
			Limit: limitInt,
			Page:  pageInt,
		},
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
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, err.Error()))
		return
	}

	isValid, message := lib.ValidateForm(data)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, message))
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
		json.NewEncoder(w).Encode(lib.NewErrorResponse(404, "User not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewDataResponse(200, "OK", data),
	)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, err.Error()))
		return
	}
	isValid, message := validateUserForm(data)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, message))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewDataResponse(201, "OK", data),
	)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Invalid Id"))
		return
	}

	if intId == 0 || intId > 100 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(404, "User not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewDataResponse(200, "DELETED", nil),
	)
}

func handleFindOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Invalid Id"))
	}
	if intId == 0 || intId > 100 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(404, "User not found"))
		return
	}
	userList := GenerateUsers(2)
	user := userList[0]
	user.Id = intId
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lib.NewDataResponse(200, "OK", user))
}

func insertMockData(database *sql.DB, logger *lib.Logger) {
	countQuery := `SELECT COUNT(*) FROM users;`
	var count int
	tx, err := database.Begin()
	if err != nil {
		logger.Error("Error creating transaction" + err.Error())
		return
	}
	tx.QueryRow(countQuery).Scan(&count)
	if count > 99 {
		logger.Info("User data already inserted Skipping")
		tx.Rollback()
		return
	}
	userList := GenerateUsers(100)
	vals := []interface{}{}
	sqlStr := `INSERT INTO users (first_name, last_name, phone, email, age, image, password) VALUES `
	for idx, user := range userList {
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", idx*7+1, idx*7+2, idx*7+3, idx*7+4, idx*7+5, idx*7+6, idx*7+7)
		vals = append(vals, user.FirstName, user.LastName, user.Phone, user.Email, user.Age, user.Image, user.Password)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += ";"
	fmt.Println(sqlStr)
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		logger.Error("Error preparing statement" + err.Error())
		return
	}
	res, err := stmt.Exec(vals...)
	if err != nil {
		tx.Rollback()
	}
	fmt.Println(res.LastInsertId())
	tx.Commit()
	logger.Info("Data inserted successfully for User")
}

func initializeMockData(database *sql.DB, logger *lib.Logger) {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			phone VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			age INTEGER,
			image VARCHAR(255),
			password VARCHAR(255) NOT NULL
		);`
	_, err := database.Exec(createTableQuery)
	if err != nil {
		logger.Info("Error creating table" + err.Error())
		log.Fatal(err)
	}
	logger.Info("Table created successfully for User")
}

func InitUserRouter(mux *mux.Router, database *sql.DB, logger *lib.Logger) {
	router := mux.PathPrefix("/api/users").Subrouter()
	initializeMockData(database, logger)
	insertMockData(database, logger)
	router.HandleFunc("", handleGet).Methods("GET")
	router.HandleFunc("", handlePost).Methods("POST")
	router.HandleFunc("/{id}", handleFindOne).Methods("GET")
	router.HandleFunc("/{id}", handlePut).Methods("PUT")
	router.HandleFunc("/{id}", handleDelete).Methods("DELETE")
}
