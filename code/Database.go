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
	name  string
	email string
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

	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
