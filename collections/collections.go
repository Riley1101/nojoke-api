package collections

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nojoke/auth"
	"nojoke/lib"

	"github.com/gorilla/mux"
)

type CollectionType struct {
	Id          int    `json:"id"`
	CreateAt    string `json:"create_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func initializeDatabase(database *sql.DB, logger *lib.Logger) {
	_, err := database.Exec(CreateProductCollectionTableQuery)
	if err != nil {
		logger.Error("Error creating collection table: " + err.Error())
	}
	logger.Info("Created collection table !")
}

func handleGet(
	w http.ResponseWriter, r *http.Request,
	database *sql.DB,
	logger *lib.Logger,
	admin *auth.Admin) {

	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	total := r.URL.Query().Get("total")

	limitInt, pageInt, totalInt, error := lib.PaginationParams(limit, page, total)

	fmt.Println(limitInt, pageInt, totalInt)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Invalid query params"))
		return
	}
	tx, error := database.Begin()
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Error creating transition"))
		return
	}
	var count int
	error = tx.QueryRow(CountProductCollectionsQuery).Scan(&count)
	if error != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Error getting count"))
		return
	}
	rows, error := tx.Query(GetProductCollectionsQuery)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Error getting collections"))
		return
	}
	collections := make([]ProductCollection, 0)
	for rows.Next() {
		collection := ProductCollection{}
		var products []uint8
		error = rows.Scan(&collection.Id, &collection.CreateAt, &products, &collection.UserId)
		fmt.Println(error)
		if error != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(lib.NewErrorResponse(400, "Error getting collections"))
			return
		}
		collection.Products = products
		collections = append(collections, collection)
	}
	defer rows.Close()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lib.NewDataResponse(200, "Success", collections))

}

func InitCollectionRouter(mux *mux.Router, database *sql.DB, logger *lib.Logger) {
	router := mux.PathPrefix("/api/collections").Subrouter()
	initializeDatabase(database, logger)
	router.Handle("", auth.Authenticated(func(w http.ResponseWriter, r *http.Request, a *auth.Admin) {
		handleGet(w, r, database, logger, a)
	})).Methods("GET")
}
