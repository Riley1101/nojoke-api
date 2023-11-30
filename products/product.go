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
	Id            int64   `json:"id"`
	Name          string  `json:"name" validate:"required"`
	Price         int     `json:"price" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	Discount      float32 `json:"discount"`
	Rating        float32 `json:"rating"`
	Stock         int     `json:"stock"`
	Brand         string  `json:"brand" validate:"required"`
	Category_id   int     `json:"category"`
	Thumbnail     string  `json:"thumbnail"`
	Image         string  `json:"image"`
	Collection_id int64   `json:"collection_id"`
}

func GenerateProducts(limit int) []Product {
	productList := []Product{}
	for i := 0; i < limit; i++ {
		Product := Product{}
		Product.Id = int64(i)
		Product.Name = faker.FirstName()
		Product.Price = rand.Intn(1000000) + 1000000
		Product.Description = faker.Paragraph()
		Product.Discount = rand.Float32()
		Product.Rating = rand.Float32() * 5
		Product.Stock = rand.Intn(100)
		Product.Brand = faker.FirstName()
		Product.Category_id = rand.Intn(100)
		Product.Thumbnail = faker.URL()
		Product.Image = faker.URL()
		productList = append(productList, Product)
	}
	return productList
}

func handleGet(w http.ResponseWriter, r *http.Request, admin *auth.Admin, database *sql.DB, logger *lib.Logger) {

	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	total := r.URL.Query().Get("total")

	limitInt, pageInt, totalInt, error := lib.PaginationParams(limit, page, total)

	query := GetProductsByCollectionQuery
	collectionId := 1
	tx, error := database.Begin()

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Error creating transition"))
		return
	}

	rows, error := tx.Query(query, collectionId)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Error getting products"))
		return
	}
	productList := []Product{}
	for rows.Next() {
		val, error := rows.Columns()
		fmt.Println(val)
		fmt.Println(error)
		product := Product{}
		error = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Discount,
			&product.Rating,
			&product.Stock,
			&product.Brand,
			&product.Category_id,
			&product.Thumbnail,
			&product.Image,
			&product.Collection_id,
		)
		productList = append(productList, product)
	}

	w.Header().Set("Content-Type", "application/json")

	response := lib.DataResponse{
		Status:  200,
		Message: "OK",
		Data:    productList,
		Pagination: lib.Pagination{
			Total: totalInt,
			Limit: limitInt,
			Page:  pageInt,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePost(w http.ResponseWriter, r *http.Request, admin *auth.Admin) {
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

func handlePut(w http.ResponseWriter, r *http.Request, admin *auth.Admin) {
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

func handleDelete(w http.ResponseWriter, r *http.Request, admin *auth.Admin) {

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
	createTableQuery := CreateProductTableQuery
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
	router.Handle("", auth.Authenticated(func(w http.ResponseWriter, r *http.Request, a *auth.Admin) {
		handleGet(w, r, a, database, logger)
	})).Methods("GET")
	router.Handle("", auth.Authenticated(handlePost)).Methods("POST")
	router.Handle("/{id}", auth.Authenticated(handlePut)).Methods("PUT")
	router.Handle("/{id}", auth.Authenticated(handleDelete)).Methods("DELETE")
}
