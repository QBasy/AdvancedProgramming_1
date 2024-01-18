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
	http.HandleFunc("/register", getMain)
	http.HandleFunc("/login", getLogin)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/404/", get404)
	http.HandleFunc("/404", get404)
	http.HandleFunc("/notfound", get404)
	http.HandleFunc("/notfound/", get404)
	http.HandleFunc("/reg", postRegister)
	http.HandleFunc("/log", postLogin)
	http.HandleFunc("/success", getSuccess)
	http.HandleFunc("/registered", getRegistered)

	fmt.Printf("Server listening on port %s...\n", port)
	http.ListenAndServe(port, nil)
}

func get404(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/404.html")
}

func getRegistered(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/registered.html")
}

func getSuccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/success.html")
}

func getMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/register.html")
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/login.html")
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := Database{}
	db.getUser(user)

	response := Response{
		Status:  "success",
		Message: "Login successful",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := Database{Name: user.Name, Password: user.Password}
	db.addUser()

	response := Response{
		Status:  "success",
		Message: "Registration successful",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	fmt.Println("User Registered")
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
