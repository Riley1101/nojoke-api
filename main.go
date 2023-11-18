package main

import (
	"fmt"
	"net/http"
	"nojoke/lib"
	"nojoke/routes"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "1337"
	}

	mux := http.NewServeMux()

	loggerMux := lib.NewLogger(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/users", routes.UserHandler)

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, loggerMux)

}
