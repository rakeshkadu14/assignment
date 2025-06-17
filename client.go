package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go <csv-filepath> <YYYY-MM>")
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/revenue", analyzeMonth).Methods("GET")

	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
