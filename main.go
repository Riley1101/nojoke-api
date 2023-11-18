package main

import (
	"fmt"
	"net/http"
	"nojoke/lib"
	"nojoke/routes"
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

	routes.InitUserRouter(r, db)

	fmt.Println("Server running on port", port)

	http.ListenAndServe(":"+port, loggerMux)

}
