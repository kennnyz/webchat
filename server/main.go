package main

import (
	"database/sql"
	"fmt"
	"github.com/kennnyz/webchat/handlers"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	// Подключение к базе данных
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=password dbname=notifications sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Создание таблицы user_actions, если она не существует
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_actions (
		ID SERIAL PRIMARY KEY,
		Username TEXT,
		Message TEXT,
		Time TIMESTAMP
	)`)

	if err != nil {
		log.Fatal("Failed to create user_actions table:", err)
	}
	routes := routes()

	log.Println("Starting channel listener")

	go handlers.ListenToWsChannel(db) // Передача подключения к базе данных в функцию ListenToWsChannel

	fmt.Println("Server is running on port 8080")
	_ = http.ListenAndServe(":8080", routes)
}
