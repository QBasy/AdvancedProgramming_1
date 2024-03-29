package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

const (
	logFilePath   = "app.log"
	logTimeFormat = "2006-01-02 15:04:05"
)

var logger *log.Logger

var limiter = rate.NewLimiter(1, 1)

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

type Videos struct {
	gorm.Model
	Name      string `json:"title"`
	Author    string `json:"author"`
	Likes     string `json:"likes"`
	Views     string `json:"views"`
	Comments  string `json:"comments"`
	Date      string `json:"date"`
	ImagePath string `json:"imagePath"`
	VideoPath string `json:"videoPath"`
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/"))))

	http.HandleFunc("/", getMain)
	http.HandleFunc("/index", getIndex)
	http.HandleFunc("/log", getLogin)
	http.HandleFunc("/register", getRegister)
	http.HandleFunc("/createUser", postRegister)
	http.HandleFunc("/notfound", get404)
	http.HandleFunc("/notfound/", get404)
	http.HandleFunc("/login", getLogin)
	http.HandleFunc("/postLogin", postLogin)
	http.HandleFunc("/profile", getProfile)
	http.HandleFunc("/registered", getRegistered)
	http.HandleFunc("/forgot", getForget)
	http.HandleFunc("/filter", getFilter)
	http.HandleFunc("/postFilter", postFilter)
	http.HandleFunc("/addVideo", getAddVideo)
	http.HandleFunc("/addVideoByUser", addVideoByUser)

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/index.html")
}

func getAddVideo(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/videoadder.html")
}

func get404(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/404.html")
}

func getRegistered(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/registered.html")
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/profile.html")
}

func getMain(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/starter.html")
}

func getForget(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/forgot.html")
}

func getRegister(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/register.html")
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/login.html")
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

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
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
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

func addVideoByUser(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var video Videos

	err := json.NewDecoder(r.Body).Decode(&video)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		fmt.Println("Error decoding JSON:", err)
		return
	}

	result := db.Create(&video)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add video")
		return
	}
	fmt.Println("success")
	http.Redirect(w, r, "/filter", http.StatusFound)
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

func getFilter(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	http.ServeFile(w, r, "frontend/filter.html")
}
func postFilter(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var filter struct {
		title  string `json:"title"`
		author string `json:"author"`
	}

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Println(filter.author + " and " + filter.author)
	query := "SELECT * FROM videos"

	if filter.title != "" {
		query += " WHERE name LIKE '%" + filter.title + "%'"
		fmt.Println(query)
	}

	if filter.author != "" {
		if filter.title != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		fmt.Println(query)
		query += " author LIKE '%" + filter.author + "%'"
	}

	query += " ORDER BY views DESC"

	fmt.Println(query)
	rows, err := db.Raw(query).Rows()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var fetchedData []Videos
	for rows.Next() {
		var video Videos

		fetchedData = append(fetchedData, video)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fetchedData)
}
