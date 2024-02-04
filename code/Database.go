package main

/*err := rows.Scan(&video.Name, &video.Author, &video.Likes, &video.Views, &video.Comments, &video.Date, &video.ImagePath, &video.VideoPath, &video.CreatedAt, &video.UpdatedAt, &video.DeletedAt, &video.ID)
imagePath := ""
videoPath := ""

video.VideoPath = videoPath
video.ImagePath = imagePath

err := rows.Scan(&video.Name, &video.Author, &video.Likes, &video.Views, &video.Comments, &video.Date, &video.ImagePath, &video.VideoPath, &video.CreatedAt, &video.UpdatedAt, &video.DeletedAt, &video.ID)

if err != nil {
	fmt.Println("ERROR", err)
	http.Error(w, "Error scanning row", http.StatusInternalServerError)
	return
}*/

/*
import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"os"
)

const connection = "host=localhost port=5432 user=postgres password=japierdole dbname=advanced sslmode=disable"

type Database struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (db Database) isInDatabase(user User) bool {
	database, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	fmt.Println("EcTb KonTakT")

	nameRows, err := database.Query("SELECT password FROM users WHERE password = $1", user.Name)
	if nameRows.Next() {
		return false
	}

	emailRows, err := database.Query("SELECT password FROM users WHERE password = $1", user.Password)
	if emailRows.Next() {
		return false
	}

	err = database.Ping()
	return true
}

func (db Database) getUser(user User) {
	query := "SELECT id, name, password FROM users WHERE name = $1 AND password = $2;"
	database, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EcTb KonTakT")

	rows, err := database.Query(query, user.Name, user.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var people []Database

	for rows.Next() {
		var user Database
		err := rows.Scan(&user.ID, &user.Name, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(people)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("file.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JSON data written, or written over old data to file.json")
}

func (db Database) addUser() {
	query := "INSERT INTO users (name, email) VALUES ($1, $2);"
	database, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EcTb KonTakT")

	if db.isInDatabase(User{db.Name, db.Password}) {
		fmt.Println("User already exists.")
		return
	}

	_, err = database.Exec(query, db.Name, db.Password)
	if err != nil {
		log.Fatal(err)
	}
}*/
