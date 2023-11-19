package products

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"nojoke/lib"

	faker "github.com/bxcodec/faker/v3"
)

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	Description string   `json:"description"`
	Discount    float32  `json:"discount"`
	Rating      float32  `json:"rating"`
	Stock       int      `json:"stock"`
	Brand       string   `json:"brand"`
	Category    Category `json:"category"`
	Thumbnail   string   `json:"thumbnail"`
	Image       string   `json:"image"`
}

func GenerateProducts(limit int) []Product {
	productList := []Product{}
	for i := 0; i < limit; i++ {
		Product := Product{}
		Product.ID = i
		Product.Name = faker.FirstName()
		Product.Price = rand.Intn(1000000) + 1000000
		Product.Description = faker.Paragraph()
		Product.Discount = rand.Float32()
		Product.Rating = rand.Float32() * 5
		Product.Stock = rand.Intn(100)
		Product.Brand = faker.FirstName()
		Product.Category = Category{
			ID:          faker.ID,
			Name:        faker.FirstName(),
			Description: faker.Paragraph(),
		}
		Product.Thumbnail = faker.URL()
		Product.Image = faker.URL()
		productList = append(productList, Product)
	}
	return productList
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
		Pagination: lib.Pagination{
			Total: totalInt,
			Limit: limitInt,
			Page:  pageInt,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
