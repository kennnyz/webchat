package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Notifier interface {
	Notify()
}

type ClientNotifier struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

func (client *ClientNotifier) Notify() {
	// SOME IMPLEMENTATION

	// ADD IN THE FUTURE MESSAGE BROKER WHERE OUR NOTIFIER WILL WRITE AND LISTEN
}

func (client *ClientNotifier) WriteActionToDB(db *sql.DB) error {
	// Подготовка SQL-запроса для вставки записи в таблицу user_actions
	stmt, err := db.Prepare("INSERT INTO user_actions (Username, Message, Time) VALUES ($1, $2, NOW())")
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement, %v", err)
	}
	defer stmt.Close()

	// Выполнение SQL-запроса с передачей параметров
	_, err = stmt.Exec(client.Username, client.Message)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)

	}

	log.Println("Successfully inserted action into user_actions table")
	return nil
}
