package routes

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"nojoke/lib"

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

func handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	v := validate.Struct(data)
	if !v.Validate() || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message := ""
		if err != nil {
			message = err.Error()
		} else {
			message = v.Errors.One()
		}
		json.NewEncoder(w).Encode(lib.NewResponse(400, message, nil))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewResponse(200, "OK", data),
	)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	case "PUT":
	case "DELETE":
	default:
		http.Error(w, "Method not allowed", 405)
	}

}

func InitUserRouter(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", UserHandler)
}
