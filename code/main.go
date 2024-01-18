package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

const port = ":8888"

var db *gorm.DB

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func init() {
	dsn := "host=localhost port=5432 user=postgres password=japierdole dbname=advanced sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", getMain)
	http.HandleFunc("/log", getLogin)
	http.HandleFunc("/register", postRegister)
	http.HandleFunc("/notfound", get404)
	http.HandleFunc("/notfound/", get404)
	http.HandleFunc("/login", postLogin)
	http.HandleFunc("/success", getSuccess)
	http.HandleFunc("/registered", getRegistered)

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return
	}
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
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	db.First(&user, "name = ? AND password = ?", user.Name, user.Password)

	if user.ID == 0 {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	response := Response{
		Status:  "success",
		Message: "Login successful",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if !isUserInDatabase(user) {
		handle404(w, r)
		return
	}

	db.Create(&user)

	http.Redirect(w, r, "/registered", http.StatusFound)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/404.html")
}

func isUserInDatabase(user User) bool {
	var existingUser User
	result := db.First(&existingUser, "name = ? OR password = ?", user.Name, user.Password)

	return result.RowsAffected == 0
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	response := Response{
		Status:  fmt.Sprint(code),
		Message: message,
	}
	respondWithJSON(w, code, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}
