package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"workoutstudy_chatting/model"
)

type UserRepository interface {
	SaveUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(userID int) error
	GetUserByID(userID int) (*model.User, error)
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

var _ UserRepository = &UserRepositoryImpl{}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (repo *UserRepositoryImpl) SaveUser(user *model.User) (*model.User, error) {
	query := `INSERT INTO "user" (id, nickname, state, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := repo.DB.QueryRow(query, user.ID, user.Nickname, user.State, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		return nil, fmt.Errorf("error saving user: %w", err)
	}
	return user, nil
}

func (repo *UserRepositoryImpl) UpdateUser(user *model.User) (*model.User, error) {
	query := `UPDATE user SET nickname = $2, state = $3, updated_at = $4 WHERE id = $1 RETURNING id`

	// 쿼리 실행
	err := repo.DB.QueryRow(query, user.ID, user.Nickname, user.State, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, fmt.Errorf("error updating user: %w", err)
	}
	return user, nil
}

func (repo *UserRepositoryImpl) DeleteUser(userID int) error {
	query := `DELETE FROM user WHERE id = $1`

	// 쿼리 실행
	_, err := repo.DB.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (repo *UserRepositoryImpl) GetUserByID(userID int) (*model.User, error) {
	query := `SELECT id, nickname, state, created_at, updated_at FROM user WHERE id = $1`

	// 쿼리 실행
	user := model.User{}
	err := repo.DB.QueryRow(query, userID).Scan(&user.ID, &user.Nickname, &user.State, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found for ID: %v", userID)
			return nil, fmt.Errorf("no user found for ID: %d", userID)
		}
		log.Printf("Error querying user by ID: %v", err)
		return nil, fmt.Errorf("error querying user by ID: %w", err)
	}
	return &user, nil
}
