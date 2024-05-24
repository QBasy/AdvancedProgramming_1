package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var limiter = rate.NewLimiter(1, 10)

const port = "8888"
const logFilePath = "logs/app.log"
const messageFilePath = "./messages/messages.json"

var allowedIPs = map[string]bool{
	"10.202.2.210": true,
	"127.0.0.1":    true,
	"localhost":    true,
}

var (
	db     *gorm.DB
	logger *logrus.Logger
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var (
	upgrader  = websocket.Upgrader{}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []byte)
	mutex     = sync.Mutex{}
)

type smtpEmail struct {
	username string
	password string
	port     string
	host     string
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
	Date      string `json:"date"`
	ImagePath string `json:"imagepath"`
}

func sendEmail(to, subject, body string, email *smtpEmail, wg *sync.WaitGroup) {
	defer wg.Done()

	from := email.username
	toAddr := []string{to}
	message := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", email.username, email.password, email.host)

	err := smtp.SendMail(email.host+":"+email.port, auth, email.username, toAddr, message)
	if err != nil {
		logrus.WithError(err).Errorf("Error sending email to %s", to)
		return
	}

	logrus.Infof("Email sent successfully to %s", to)
}

func init() {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	logger = logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{})

	const dbURL = "postgres://qbasy:ZDDwUSlj2L6TPRjXMsW8Y3gS2uPRNEkJ@dpg-cni5e2n79t8c73bprang-a.frankfurt-postgres.render.com/advanced"
	var err1 error
	db, err1 = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err1 != nil {
		logger.Fatalf("Error connecting to database: %v", err1)
	}

	err2 := db.AutoMigrate(&User{})
	if err2 != nil {
		logger.Fatalf("Error migrating database: %v", err2)
	}

	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)
}

func createTablesIfNotExists(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&Photos{})
	if err != nil {
		return
	}
}

func main() {
	initLogger()
	initDB()

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/ws", handleWebSocket)

	go broadcaster()

	srv := &http.Server{Addr: ":" + port}

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

func initLogger() {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	logger = logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func initDB() {
	const dbURL = "postgres://qbasy:ZDDwUSlj2L6TPRjXMsW8Y3gS2uPRNEkJ@dpg-cni5e2n79t8c73bprang-a.frankfurt-postgres.render.com/advanced"
	var err1 error
	db, err1 = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err1 != nil {
		logger.Fatalf("Error connecting to database: %v", err1)
	}

	createTablesIfNotExists(db)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- message
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				log.Println("Error writing message:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/" {
		r.URL.Path += "starter"
		http.Redirect(w, r, "/starter", 500)
	}

	switch r.URL.Path {
	case "/admin":
		if !isAdminAllowed(ip) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		serveAdmin(w, r)
	case "/index":
		serveIndex(w, r)
	case "/send-message":
		sendMessageHandler(w, r)

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

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	log.Printf("Received message from admin: %s", message.Message)

	if err := saveMessageToFile(message); err != nil {
		http.Error(w, "Error saving message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Message received and saved"})
}

func saveMessageToFile(message Message) error {
	file, err := os.OpenFile("messages/"+message.Username+".json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(message); err != nil {
		return err
	}
	defer file.Close()

	return nil
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
	http.ServeFile(w, r, "frontend/index.html")
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

func serveAdmin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/admin.html")
}

func isAdminAllowed(ip string) bool {
	fmt.Println(ip)
	allowedIPs := map[string]bool{
		"10.202.2.210": true,
		"10.202.6.198": true,
		"127.0.0.1":    true,
		"localhost":    true,
		"::1":          true,
	}
	return allowedIPs[ip]
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

	if isUserInDatabase(user) {
		respondWithError(w, http.StatusConflict, "User already exists")
		return
	}

	db.Create(&user)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send welcome email")
		return
	}

	response := Response{
		Status:  "success",
		Message: "User created successfully",
	}
	logrus.Println("User created")
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

	// Query the database to find the user
	result := db.First(&user, "name = ? AND password = ?", user.Name, user.Password)
	if result.RowsAffected == 0 {
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

	dbPhoto := Photos{
		Name:      photo.Name,
		Author:    photo.Author,
		Likes:     photo.Likes,
		Date:      photo.Date,
		ImagePath: photo.ImagePath,
	}

	result := db.Create(&dbPhoto)
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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var fetchedData []Photos
	for rows.Next() {
		var photo Photos
		err := db.ScanRows(rows, &photo)
		if err != nil {
			return
		}
		fetchedData = append(fetchedData, photo)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fetchedData)
	if err != nil {
		return
	}
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
