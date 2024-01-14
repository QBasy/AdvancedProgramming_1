package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const port = ":8888"
const databaseStr = "user=postgres dbname=advanced sslmode=disable"

func main() {
	http.HandleFunc("/", getPostRequest)
	fmt.Printf("Server listening on port %s...\n", port)
	http.ListenAndServe(port, nil)
}

func getResponse(w http.ResponseWriter, r *http.Request) {
}

func getMain(w http.ResponseWriter, r *http.Request) {
	w.Header()
}

func getPostRequest(w http.ResponseWriter, r *http.Request) {
	var request Request

	fmt.Printf("message: %s\n", request.Message)

	response := Response{
		Status:  "success",
		Message: "Data successfully received",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
