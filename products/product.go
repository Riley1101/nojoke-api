package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"nojoke/auth"
	"nojoke/lib"
	"strconv"

	"github.com/gorilla/mux"

	faker "github.com/bxcodec/faker/v3"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	Id          int      `json:"id"`
	Name        string   `json:"name" validate:"required"`
	Price       int      `json:"price" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Discount    float32  `json:"discount"`
	Rating      float32  `json:"rating"`
	Stock       int      `json:"stock"`
	Brand       string   `json:"brand" validate:"required"`
	Category    Category `json:"category"`
	Thumbnail   string   `json:"thumbnail"`
	Image       string   `json:"image"`
}

func GenerateProducts(limit int) []Product {
	productList := []Product{}
	for i := 0; i < limit; i++ {
		Product := Product{}
		Product.Id = i
		Product.Name = faker.FirstName()
		Product.Price = rand.Intn(1000000) + 1000000
		Product.Description = faker.Paragraph()
		Product.Discount = rand.Float32()
		Product.Rating = rand.Float32() * 5
		Product.Stock = rand.Intn(100)
		Product.Brand = faker.FirstName()
		Product.Category = Category{
			ID:          rand.Intn(100),
			Name:        faker.FirstName(),
			Description: faker.Paragraph(),
		}
		Product.Thumbnail = faker.URL()
		Product.Image = faker.URL()
		productList = append(productList, Product)
	}
	return productList
}

func handleGet(w http.ResponseWriter, r *http.Request, admin *auth.Admin) {

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

	productList := GenerateProducts(totalInt)
	paginatedProducts := lib.PaginateData(productList, limitInt, pageInt, totalInt)

	response := lib.DataResponse{
		Status:  200,
		Message: "OK",
		Data:    paginatedProducts,
		Pagination: lib.Pagination{
			Total: totalInt,
			Limit: limitInt,
			Page:  pageInt,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data := Product{}
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewDataResponse(201, "OK", data),
	)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	pathParams := mux.Vars(r)
	id := pathParams["id"]
	data := Product{}
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
	userList := GenerateProducts(100)
	found := false
	for i, product := range userList {
		if product.Id == data.Id || fmt.Sprint(product.Id) == id {
			userList[i] = data
			found = true
			break
		}
	}
	if !found || data.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(404, "Product not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		lib.NewDataResponse(200, "OK", data),
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

func initializeDatabase(database *sql.DB, logger *lib.Logger) {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price INT NOT NULL,
			description TEXT NOT NULL,
			discount FLOAT,
			rating FLOAT,
			stock INT NOT NULL,
			brand VARCHAR(255) NOT NULL,
			category_id INT,
			thumbnail VARCHAR(255),
			image VARCHAR(255)
		);`
	_, err := database.Exec(createTableQuery)
	if err != nil {
		logger.Error("Error creating table" + err.Error())
		return
	}
	logger.Info("Table created successfully for Products")
}

func InitProductRouter(mux *mux.Router, database *sql.DB, logger *lib.Logger) {
	initializeDatabase(database, logger)
	router := mux.PathPrefix("/api/products").Subrouter()
	router.Handle("", auth.Authenticated(handleGet)).Methods("GET")
	router.HandleFunc("", handlePost).Methods("POST")
	router.HandleFunc("/{id}", handlePut).Methods("PUT")
	router.HandleFunc("/{id}", handleDelete).Methods("DELETE")
}
