package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "postgresql-chatting"
	port     = 5432
	user     = "chatting"
	password = "chatting"
	dbname   = "chatting-db"
)

func InitializeDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	createTables := []string{
		`CREATE TABLE IF NOT EXISTS fit_group (
			id SERIAL PRIMARY KEY,
			fit_group_name VARCHAR(30),
			max_fit_mate INTEGER,
			created_at TIMESTAMP NOT NULL,
			created_by VARCHAR(30),
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(30)
		)`,
		`CREATE TABLE IF NOT EXISTS fit_mate (
			id SERIAL PRIMARY KEY,
			fit_group_id INT REFERENCES fit_group(id),
			nickname VARCHAR(10),
			state BOOLEAN,
			created_at TIMESTAMP NOT NULL,
			created_by VARCHAR(30),
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(30)
		)`,
	}

	for _, query := range createTables {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Failed to execute query: %v, error: %v", query, err)
		}
	}

	fmt.Println("Database initialized successfully")
}
