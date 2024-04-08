package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"workoutstudy_chatting/model"
)

type FitGroupRepository interface {
	GetFitGroupByID(id int) (*model.FitGroup, error)
}

type FitGroupRepositoryImpl struct {
	DB *sql.DB
}

// 훈기 tip : 인터페이스 메소드 슬라이스 중 구현안되거 있으면 에러 띄워줌
var _ FitGroupRepository = &FitGroupRepositoryImpl{}

func NewFitGroupRepository(db *sql.DB) FitGroupRepository {
	return &FitGroupRepositoryImpl{DB: db}
}

func (repo *FitGroupRepositoryImpl) GetFitGroupByID(id int) (*model.FitGroup, error) {
	query := `SELECT id, fit_group_name, max_fit_mate, created_at, created_by, updated_at, updated_by FROM fit_group WHERE id = $1`

	log.Printf("Repository layer: Executing query for FitGroupID: %d", id)
	fitGroup := model.FitGroup{}
	err := repo.DB.QueryRow(query, id).Scan(&fitGroup.ID, &fitGroup.FitGroupName, &fitGroup.MaxFitMate, &fitGroup.CreatedAt, &fitGroup.CreatedBy, &fitGroup.UpdatedAt, &fitGroup.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Repository layer: No fit_group found for ID: %v", id)
			return nil, fmt.Errorf("no fit_group found for ID: %d", id)
		}
		log.Printf("Repository layer: Error querying fit_group by ID: %v", err)
		return nil, err
	}

	return &fitGroup, nil
}
