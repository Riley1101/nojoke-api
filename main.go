package main

import (
	"fmt"
	"net/http"
	auth "nojoke/auth"
	"nojoke/collections"
	"nojoke/lib"
	product "nojoke/products"
	users "nojoke/users"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	db := lib.ConnectDB()
	if port == "" {
		port = "1337"
	}

	r := mux.NewRouter()
	loggerMux := lib.NewLogger(r)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	auth.InitAuthRouter(r, db, loggerMux)

	users.InitUserRouter(r, db, loggerMux)

	product.InitProductRouter(r, db, loggerMux)

	collections.InitCollectionRouter(r, db, loggerMux)

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, loggerMux)

}
