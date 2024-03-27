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

var DB *sql.DB

func InitializeDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	// 데이터베이스 연결 풀 설정
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(0)

	// 데이터베이스 연결 테스트
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Database connection pool initialized successfully")

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
			username VARCHAR(20),
			nickname VARCHAR(10),
			state BOOLEAN,
			created_at TIMESTAMP NOT NULL,
			created_by VARCHAR(30),
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(30)
		)`,
	}

	for _, query := range createTables {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("Failed to execute query: %v, error: %v", query, err)
		}
	}

	fmt.Println("Database initialized successfully")

	// 더미 데이터 삽입
	insertDummyData := []string{
		`INSERT INTO fit_group (fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by)
		VALUES ('운터디', 20, NOW(), '서경원', NOW(), '서경원');`,
		`INSERT INTO fit_group (fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by)
		VALUES ('축구의신', 11, NOW(), '손흥민', NOW(), '손흥민');`,
		`INSERT INTO fit_mate (fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by)
        VALUES (1, '서경원', '경원이', TRUE, NOW(), '서경원', NOW(), '서경원') ON CONFLICT (id) DO NOTHING;`,
		`INSERT INTO fit_mate (fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by)
        VALUES (2, '손흥민', '대흥민', TRUE, NOW(), '손흥민', NOW(), '손흥민') ON CONFLICT (id) DO NOTHING;`,
	}

	// 더미 데이터 삽입 실행
	for _, query := range insertDummyData {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("Failed to insert dummy data: %v, error: %v", query, err)
		}
	}

	// 나머지 8명의 사용자에 대한 더미 데이터 삽입
	for i := 3; i <= 10; i++ {
		query := fmt.Sprintf(`INSERT INTO fit_mate (fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by)
        VALUES (%d, '테스트유저%d', '테스트유저별명%d', TRUE, NOW(), '테스트유저%d', NOW(), '테스트유저%d') ON CONFLICT (id) DO NOTHING;`, i%2+1, i, i, i, i)
		// fit_group_id를 교대로 1과 2로 설정합니다. i가 홀수일 때는 1, 짝수일 때는 2로.

		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("Failed to insert dummy data for fit_mate: %v, error: %v", query, err)
		}
	}
	fmt.Println("Database initialized successfully with dummy data")
}
