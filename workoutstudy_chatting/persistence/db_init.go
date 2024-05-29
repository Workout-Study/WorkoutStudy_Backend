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

func InitializeDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

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
		`CREATE TABLE IF NOT EXISTS "user" (
			id INTEGER PRIMARY KEY,
			nickname VARCHAR(10) NOT NULL,
			state BOOLEAN DEFAULT false NOT NULL,
			created_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			created_by VARCHAR(30),
			updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			updated_by VARCHAR(30)
		)`,
		`CREATE TABLE IF NOT EXISTS fit_group (
			id INTEGER PRIMARY KEY,
			fit_leader_user_id INTERGER REFERENCES user(id) NOT NULL,
			fit_group_name VARCHAR(30),
			category INTEGER NOT NULL,
			cycle INTEGER NOT NULL,
			frequency INTEGER NOT NULL,
			present_fit_mate_count INTEGER NOT NULL,
			max_fit_mate INTEGER NOT NULL,
			state BOOLEAN DEFAULT false NOT NULL,
			created_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			created_by VARCHAR(30),
			updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			updated_by VARCHAR(30)
		)`,
		`CREATE TABLE IF NOT EXISTS fit_mate (
			id INTEGER PRIMARY KEY,
			fit_group_id INTEGER REFERENCES fit_group(id) NOT NULL,
			state BOOLEAN DEFAULT false NOT NULL,
			create_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			create_by VARCHAR(30),
			updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			updated_by VARCHAR(30)
		)`,
		`CREATE TABLE IF NOT EXISTS message (
			message_id UUID PRIMARY KEY,
			user_id INTEGER REFERENCES user(id) NOT NULL,
			fit_group_id INT REFERENCES fit_group(id) NOT NULL,
			message TEXT NOT NULL,
			message_time TIMESTAMP(6),
			message_type VARCHAR(8) CHECK (message_type IN ('CHATTING', 'TICKET')),
			created_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			created_by VARCHAR(30),
			updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL,
			updated_by VARCHAR(30)
		)`,
	}

	for _, query := range createTables {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Failed to execute query: %v, error: %v", query, err)
		}
	}

	fmt.Println("Database initialized successfully")

	// 더미 데이터 삽입
	// insertDummyData := []string{
	// 	`INSERT INTO fit_group (fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by)
	// 	VALUES ('운터디', 20, NOW(), '서경원', NOW(), '서경원');`,
	// 	`INSERT INTO fit_group (fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by)
	// 	VALUES ('축구의신', 11, NOW(), '손흥민', NOW(), '손흥민');`,
	// 	`INSERT INTO fit_mate (fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by)
	//     VALUES (1, '서경원', '경원이', TRUE, NOW(), '서경원', NOW(), '서경원') ON CONFLICT (id) DO NOTHING;`,
	// 	`INSERT INTO fit_mate (fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by)
	//     VALUES (2, '손흥민', '대흥민', TRUE, NOW(), '손흥민', NOW(), '손흥민') ON CONFLICT (id) DO NOTHING;`,
	// 	`INSERT INTO fit_group_mate (fit_group_id, fit_mate_id) VALUES (1, 1);`,
	// 	`INSERT INTO fit_group_mate (fit_group_id, fit_mate_id) VALUES (2, 1);`,
	// }

	// 더미 데이터 삽입 실행
	// for _, query := range insertDummyData {
	// 	if _, err := DB.Exec(query); err != nil {
	// 		log.Fatalf("Failed to insert dummy data: %v, error: %v", query, err)
	// 	}
	// }

	// 나머지 8명의 사용자에 대한 더미 데이터 삽입
	// for i := 3; i <= 10; i++ {
	// 	query := fmt.Sprintf(`INSERT INTO fit_mate (fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by)
	//     VALUES (%d, '테스트유저%d', '테스트유저별명%d', TRUE, NOW(), '테스트유저%d', NOW(), '테스트유저%d') ON CONFLICT (id) DO NOTHING;`, i%2+1, i, i, i, i)
	// fit_group_id를 교대로 1과 2로 설정합니다. i가 홀수일 때는 1, 짝수일 때는 2로.

	// 	_, err := DB.Exec(query)
	// 	if err != nil {
	// 		log.Fatalf("Failed to insert dummy data for fit_mate: %v, error: %v", query, err)
	// 	}
	// }
	// fmt.Println("Database initialized successfully with dummy data")

	// 메시지 더미 데이터 삽입

	// for i := 1; i <= 20; i++ { // i의 범위를 1부터 20까지로 변경
	// 	for _, fitGroupID := range []int{1, 2} {
	// 		// 메시지의 시간을 분 단위로 조정하여 더욱 현실적으로 만듭니다.
	// 		messageText := fmt.Sprintf("안녕하세요%d", i) // 메시지 텍스트 동적 생성
	// 		query := `INSERT INTO message (message_id, fit_group_id, fit_mate_id, message, message_time, message_type, created_at, created_by, updated_at, updated_by)
	//     VALUES (gen_random_uuid(), $1, 1, $2, NOW(), 'CHATTING', NOW(), '서경원', NOW(), '서경원')`
	// 		_, err := DB.Exec(query, fitGroupID, messageText)
	// 		if err != nil {
	// 			log.Fatalf("Failed to insert dummy message data: error: %v", err)
	// 		}
	// 	}
	// }

	// fmt.Println("Message dummy data inserted successfully")

	return DB
}
