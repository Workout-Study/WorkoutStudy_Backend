package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"workoutstudy_chatting/model" // 모델 패키지 경로에 맞게 수정
)

// FitMateRepository 인터페이스 선언
type FitMateRepository interface {
	GetFitGroupByMateID(fitMateID string) ([]model.FitGroup, error)
	GetFitMateByID(fitMateID string) (*model.FitMate, error)
}

type PostgresFitMateRepository struct {
	DB *sql.DB
}

// NewPostgresFitMateRepository 생성자 함수는 PostgresFitMateRepository의 새 인스턴스를 반환합니다.
func NewPostgresFitMateRepository(db *sql.DB) *PostgresFitMateRepository {
	return &PostgresFitMateRepository{DB: db}
}

func (repo *PostgresFitMateRepository) GetFitGroupByMateID(fitMateID string) ([]model.FitGroup, error) {
	query := `
	SELECT fg.id, fg.fit_group_name, fg.max_fit_mate, fg.created_at, fg.created_by, fg.updated_at, fg.updated_by
	FROM fit_group fg
	INNER JOIN fit_group_mate fgm ON fg.id = fgm.fit_group_id
	WHERE fgm.fit_mate_id = $1
	`

	rows, err := repo.DB.Query(query, fitMateID)
	if err != nil {
		log.Printf("Error retrieving fit groups by mate ID: %v", err)
		return nil, fmt.Errorf("error retrieving fit groups by mate ID: %w", err)
	}
	defer rows.Close()

	var fitGroups []model.FitGroup
	for rows.Next() {
		var fg model.FitGroup
		if err := rows.Scan(&fg.ID, &fg.FitGroupName, &fg.MaxFitMate, &fg.CreatedAt, &fg.CreatedBy, &fg.UpdatedAt, &fg.UpdatedBy); err != nil {
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
	SELECT id, fit_group_id, username, nickname, state, created_at, created_by, updated_at, updated_by
	FROM fit_mate
	WHERE id = $1
	`
	var fm model.FitMate
	err := repo.DB.QueryRow(query, fitMateID).Scan(&fm.ID, &fm.FitGroupID, &fm.Username, &fm.Nickname, &fm.State, &fm.CreatedAt, &fm.CreatedBy, &fm.UpdatedAt, &fm.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &fm, nil
}
