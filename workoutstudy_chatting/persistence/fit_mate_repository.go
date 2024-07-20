package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"workoutstudy_chatting/model" // 모델 패키지 경로에 맞게 수정
)

// FitMateRepository 인터페이스 선언
type FitMateRepository interface {
	GetFitGroupsByUserID(userID int) ([]model.FitGroup, error)
	GetFitMateByID(fitMateID string) (*model.FitMate, error)
	SaveFitMate(fitMate *model.FitMate) (*model.FitMate, error)
	DeleteFitMate(id int) error
	UpdateFitMate(fitMate *model.FitMate) (*model.FitMate, error)
	GetFitMatesIdsByFitGroupId(fitGroupId int) ([]int, error)
	CheckFitGroupExists(fitGroupID int) (bool, error)
}

type PostgresFitMateRepository struct {
	DB *sql.DB
}

// NewPostgresFitMateRepository 생성자 함수는 PostgresFitMateRepository의 새 인스턴스를 반환합니다.
func NewPostgresFitMateRepository(db *sql.DB) *PostgresFitMateRepository {
	return &PostgresFitMateRepository{DB: db}
}

func (repo *PostgresFitMateRepository) GetFitGroupsByUserID(userID int) ([]model.FitGroup, error) {
	query := `
	SELECT fg.id, fg.fit_leader_user_id, fg.fit_group_name, fg.category, fg.cycle, fg.frequency, fg.present_fit_mate_count, fg.max_fit_mate, fg.created_at, fg.created_by, fg.updated_at, fg.updated_by
	FROM fit_group fg
	INNER JOIN fit_mate fm ON fg.id = fm.fit_group_id
	WHERE fm.user_id = $1
	`

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving fit groups by user ID: %v", err)
		return nil, fmt.Errorf("error retrieving fit groups by user ID: %w", err)
	}
	defer rows.Close()

	var fitGroups []model.FitGroup
	for rows.Next() {
		var fg model.FitGroup
		if err := rows.Scan(&fg.ID, &fg.FitLeaderUserID, &fg.FitGroupName, &fg.Category, &fg.Cycle, &fg.Frequency, &fg.PresentFitMateCount, &fg.MaxFitMate, &fg.CreatedAt, &fg.CreatedBy, &fg.UpdatedAt, &fg.UpdatedBy); err != nil {
			log.Printf("Error scanning fit group: %v", err)
			continue // or return an error; depends on your error handling strategy
		}
		fitGroups = append(fitGroups, fg)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, fmt.Errorf("repo-error: iterating rows: %w", err)
	}

	return fitGroups, nil
}

func (repo *PostgresFitMateRepository) GetFitMateByID(fitMateID string) (*model.FitMate, error) {
	query := `
	SELECT id, user_id, fit_group_id, state, created_at, created_by, updated_at, updated_by
	FROM fit_mate
	WHERE id = $1
	`
	var fm model.FitMate
	err := repo.DB.QueryRow(query, fitMateID).Scan(&fm.ID, &fm.UserID, &fm.FitGroupID, &fm.State, &fm.CreatedAt, &fm.CreatedBy, &fm.UpdatedAt, &fm.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &fm, nil
}

func (repo *PostgresFitMateRepository) SaveFitMate(fitMate *model.FitMate) (*model.FitMate, error) {
	query := `INSERT INTO fit_mate (id, user_id, fit_group_id, state, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, NOW(), $5, NOW(), $5) RETURNING id`
	err := repo.DB.QueryRow(query, fitMate.ID, fitMate.UserID, fitMate.FitGroupID, fitMate.State, fitMate.CreatedBy).Scan(&fitMate.ID)
	if err != nil {
		return nil, err
	}
	return fitMate, nil
}

func (repo *PostgresFitMateRepository) DeleteFitMate(id int) error {
	query := `DELETE FROM fit_mate WHERE id = $1`
	_, err := repo.DB.Exec(query, id)
	return err
}

func (repo *PostgresFitMateRepository) UpdateFitMate(fitMate *model.FitMate) (*model.FitMate, error) {
	query := `UPDATE fit_mate SET user_id = $2, fit_group_id = $3, state = $4, updated_at = NOW(), updated_by = $5 WHERE id = $1 RETURNING id`
	err := repo.DB.QueryRow(query, fitMate.ID, fitMate.UserID, fitMate.FitGroupID, fitMate.State, fitMate.UpdatedBy).Scan(&fitMate.ID)
	if err != nil {
		return nil, err
	}
	return fitMate, nil
}

func (repo *PostgresFitMateRepository) GetFitMatesIdsByFitGroupId(fitGroupId int) ([]int, error) {
	query := `SELECT id FROM fit_mate WHERE fit_group_id = $1`
	rows, err := repo.DB.Query(query, fitGroupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fitMateIds []int
	for rows.Next() {
		var fitMateId int
		if err := rows.Scan(&fitMateId); err != nil {
			return nil, err
		}
		fitMateIds = append(fitMateIds, fitMateId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fitMateIds, nil
}

func (repo *PostgresFitMateRepository) CheckFitGroupExists(fitGroupID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fit_group WHERE id = $1)"
	var exists bool
	err := repo.DB.QueryRow(query, fitGroupID).Scan(&exists)
	return exists, err
}
