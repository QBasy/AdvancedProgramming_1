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

func main() {
	http.HandleFunc("/", getMain)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/404/", get404)
	http.HandleFunc("/404", get404)

	http.HandleFunc("/notfound", get404)
	http.HandleFunc("/notfound/", get404)

	fmt.Printf("Server listening on port %s...\n", port)
	http.ListenAndServe(port, nil)
}

func get404(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/404.html")
}

func getMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/register.html")
}

func postLogin(w http.ResponseWriter, r *http.Request) {

}

func postRegister(w http.ResponseWriter, r *http.Request) {

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
