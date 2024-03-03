package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var limiter = rate.NewLimiter(1, 10)

const port = ":8888"
const logFilePath = "logs/app.log"

var db *gorm.DB
var logger *logrus.Logger

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Photos struct {
	gorm.Model
	Name      string `json:"title"`
	Author    string `json:"author"`
	Likes     string `json:"likes"`
	Comments  string `json:"comments"`
	Date      string `json:"date"`
	ImagePath string `json:"imagepath"`
}

func init() {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	logger = logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{})

	dsn := "host=localhost port=5432 user=postgres password=japierdole dbname=advanced sslmode=disable"
	var err1 error
	db, err1 = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err1 != nil {
		logger.Fatalf("Error connecting to database: %v", err1)
	}

	err2 := db.AutoMigrate(&User{})
	if err2 != nil {
		logger.Fatalf("Error migrating database: %v", err2)
	}

	defer logFile.Close()
}

func main() {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{})

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Fatalf("Error getting database: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			logrus.Fatalf("Error closing database: %v", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Server Started")

		if r.URL.Path == "/" {
			http.FileServer(http.Dir("frontend/")).ServeHTTP(w, r)
			return
		}

		handleRequest(w, r)
	})

	srv := &http.Server{Addr: port}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		logger.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Fatalf("Error shutting down server: %v", err)
		}
		logrus.Info("Server shutdown")
	}()

	logrus.Printf("Server listening on port %s...\n", port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Error starting server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	if r.URL.Path == "/" {
		r.URL.Path += "starter"
		http.Redirect(w, r, "/starter", 500)
	}

	switch r.URL.Path {
	case "/index":
		serveIndex(w, r)
	case "/addVideo":
		serveAddVideo(w, r)
	case "/notfound":
		serveNotFound(w, r)
	case "/login":
		serveLogin(w, r)
	case "/register":
		serveRegister(w, r)
	case "/createUser":
		createUser(w, r)
	case "/postLogin":
		postLogin(w, r)
	case "/registered":
		serveRegistered(w, r)
	case "/forgot":
		serveForgot(w, r)
	case "/filter":
		serveFilter(w, r)
	case "/postFilter":
		postFilter(w, r)
	case "/starter":
		serveStarter(w, r)
	case "/addVideoByUser":
		addVideoByUser(w, r)
	default:
		serveNotFound(w, r)
	}
}

func serveStarter(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/starter.html")
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/starter.html")
}

func serveAddVideo(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/videoadder.html")
}

func serveNotFound(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/404.html")
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/login.html")
}

func serveRegister(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/register.html")
}

func createUser(w http.ResponseWriter, r *http.Request) {
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

	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Password) == "" {
		respondWithError(w, http.StatusBadRequest, "Name or Password cannot be empty")
		return
	}

	if !isUserInDatabase(user) {
		respondWithError(w, http.StatusConflict, "User already exists")
		return
	}

	db.Create(&user)

	response := Response{
		Status:  "success",
		Message: "User created successfully",
	}
	respondWithJSON(w, http.StatusCreated, response)
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

func serveRegistered(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/registered.html")
}

func serveForgot(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	http.ServeFile(w, r, "frontend/forgot.html")
}

func serveFilter(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	http.ServeFile(w, r, "frontend/filter.html")
}

func addVideoByUser(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var photo Photos

	err := json.NewDecoder(r.Body).Decode(&photo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		logger.Println("Error decoding JSON:", err)
		return
	}

	dbVideo := Photos{
		Name:      photo.Name,
		Author:    photo.Author,
		Likes:     photo.Likes,
		Comments:  photo.Comments,
		Date:      photo.Date,
		ImagePath: photo.ImagePath,
	}

	result := db.Create(&dbVideo)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add Photo")
		return
	}
	logger.WithFields(logrus.Fields{
		"action": "add_photo",
		"video":  photo.Name,
	}).Info("Photo added successfully")
	http.Redirect(w, r, "/filter", http.StatusFound)
}

func postFilter(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 12
	}

	var filter struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	err = json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	query := "SELECT * FROM videos"

	if filter.Title != "" {
		query += " WHERE name LIKE '%" + filter.Title + "%'"
	}

	if filter.Author != "" {
		if filter.Title != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " author LIKE '%" + filter.Author + "%'"
	}

	// sort
	sortOrder := r.URL.Query().Get("sortOrder")

	switch strings.ToLower(sortOrder) {
	case "1":
		query += " ORDER BY date ASC"
		break
	case "2":
		query += " ORDER BY date DESC"
		break
	default:
		query += " ORDER BY date DESC"
	}

	rows, err := db.Raw(query).Rows()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var fetchedData []Photos
	for rows.Next() {
		var photo Photos
		db.ScanRows(rows, &photo)
		fetchedData = append(fetchedData, photo)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fetchedData)
}

func isUserInDatabase(user User) bool {
	var existingUser User
	result := db.First(&existingUser, "name = ?", user.Name)

	return result.RowsAffected != 0
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
		logger.Printf("Error writing response: %v\n", err)
	}
}
