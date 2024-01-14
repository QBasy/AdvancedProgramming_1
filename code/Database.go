package main

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
	name  string
	email string
}

type User struct {
	name  string `json:"name"`
	email string `json:"email"`
}

func (db Database) isInDatabase(user User) bool {
	database, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EcTb KonTakT")

	name, err := database.Query("SELECT name FROM users WHERE name = $1", user.name)

	if name != nil {
		return false
	}

	email, err := database.Query("SELECT email FROM users WHERE email = $1", user.email)

	if email != nil {
		return false
	}

	err = database.Ping()
	return true
}
func (db Database) getUser(user User) {
	query := "SELECT users.id, users.name, users.email FROM users WHERE users.name = " + user.name + " AND users.email = " + user.email + ";"
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

	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var people []Database

	for rows.Next() {
		var user Database
		err := rows.Scan(&user.name, &user.email)
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
	query := "INSERT INTO users (users.name, users.email) VALUES (" + db.name + ", " + db.email + ");"
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

	db.isInDatabase(User{db.name, db.email})

	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
